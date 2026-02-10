package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LotteryRepository struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewLotteryRepository(db *mongo.Database, rdb *redis.Client) *LotteryRepository {
	collection := db.Collection("lotteries")

	// Create indexes
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "number", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "reserved_until", Value: 1}},
			// TTL index to automatically clear expired reservations (optional, but good for cleanup)
		},
	}
	collection.Indexes().CreateMany(context.Background(), indexModels)

	return &LotteryRepository{
		collection: collection,
		redis:      rdb,
	}
}

func (r *LotteryRepository) patternToRegex(pattern string) string {
	regex := strings.ReplaceAll(pattern, "*", ".")
	return fmt.Sprintf("^%s$", regex)
}

func (r *LotteryRepository) SearchAndReserve(ctx context.Context, pattern string, userID string, limit int) ([]domain.LotteryTicket, error) {
	// แปลงรูปแบบ wildcard (เช่น ****23) ให้เป็น Regular Expression ของ MongoDB (เช่น ^....23$)
	regexPattern := r.patternToRegex(pattern)
	redisKey := fmt.Sprintf("lottery_pattern:%s", pattern)
	now := time.Now()
	// กำหนดเวลาหมดอายุของการจอง (เช่น 5 นาที)
	reservedUntil := now.Add(5 * time.Minute)

	tickets := make([]domain.LotteryTicket, 0)

	// 1. ใช้ค่าจาก redis ก่อน
	// ใช้คำสั่ง SPopN เพื่อดึงหมายเลขออกมาแบบระบุจำนวนและรับประกันความเป็น Atomic (ใช้คนเดียวแน่นอน ดึงแล้วลบทันที)
	ticketNumbers, err := r.redis.SPopN(ctx, redisKey, int64(limit)).Result()
	if err == nil && len(ticketNumbers) > 0 {
		for _, num := range ticketNumbers {
			// เมื่อได้เลขจาก Redis แล้ว ต้องทำการอัปเดตสถานะใน MongoDB เพื่อยืนยันการจองจริง
			filter := bson.M{"number": num, "status": domain.LotteryStatusAvailable}
			update := bson.M{
				"$set": bson.M{
					"status":         domain.LotteryStatusReserved,
					"reserved_by":    userID,
					"reserved_until": reservedUntil,
					"updated_at":     now,
				},
			}
			// ใช้ FindOneAndUpdate เพื่อให้อัปเดตและดึงข้อมูลกลับมาในคำสั่งเดียวแบบ Atomic ป้องกันข้อมูลซ้ำซ้อน
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
			var doc lotteryDoc
			err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
			if err == nil {
				tickets = append(tickets, *doc.toLotteryDomain())
			}
		}
	}

	// 2. หากใน Redis มีเลขไม่พอ ให้ไปค้นหาโดยตรงจาก MongoDB
	remaining := limit - len(tickets)
	if remaining > 0 {
		for i := 0; i < remaining; i++ {
			// ค้นหาลอตเตอรี่ที่ตรงกับ Pattern และมีสถานะเป็น 'Available'
			// หรือเป็นลอตเตอรี่ที่จองไว้แล้วแต่หมดเวลาการจอง
			filter := bson.M{
				"number": bson.M{"$regex": regexPattern},
				"$or": []bson.M{
					{"status": domain.LotteryStatusAvailable},
					{
						"status":         domain.LotteryStatusReserved,
						"reserved_until": bson.M{"$lt": now},
					},
				},
			}

			update := bson.M{
				"$set": bson.M{
					"status":         domain.LotteryStatusReserved,
					"reserved_by":    userID,
					"reserved_until": reservedUntil,
					"updated_at":     now,
				},
			}

			// ทำการจองทันทีที่ค้นพบเลขที่ตรงเงื่อนไขแบบ Atomic เพื่อป้องกัน Race Condition
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
			var doc lotteryDoc
			err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					break
				}
				return nil, err
			}
			tickets = append(tickets, *doc.toLotteryDomain())
		}
	}

	// 3. เรียกทำงานเบื้องหลัง (Background Goroutine) เพื่อเติมเลขเข้า Redis สำหรับการค้นหาครั้งต่อไป
	// มันจะช่วยเพิ่มประสิทธิภาพ เมื่อมีผู้ใช้จำนวนมากค้นหาแบบเดิมพร้อมๆกัน
	go r.prefillRedis(pattern, redisKey)

	return tickets, nil
}

func (r *LotteryRepository) prefillRedis(pattern string, redisKey string) {
	// 1. ตรวจสอบปริมาณข้อมูลใน Redis Pool ของ Pattern นี้
	// หากยังมีข้อมูลเหลือมากกว่า 50 รายการ หรือเกิดข้อผิดพลาด ให้หยุดการทำงาน (เพื่อประหยัดทรัพยากร)
	count, err := r.redis.SCard(context.Background(), redisKey).Result()
	if err != nil || count > 50 {
		return
	}

	ctx := context.Background()
	regexPattern := r.patternToRegex(pattern)

	// 2. ดึงหมายเลขลอตเตอรี่ที่ว่าง (Available) จาก MongoDB เพื่อนำไปเติมใน Redis Pool
	// ค้นหาโดยใช้ Regex และจำกัดจำนวนไว้ที่ 100 รายการ เพื่อไม่ให้โหลดข้อมูลหนักเกินไป
	filter := bson.M{
		"number": bson.M{"$regex": regexPattern},
		"status": domain.LotteryStatusAvailable,
	}
	// เลือกเฉพาะฟิลด์ "number" เท่านั้น เพื่อลดปริมาณข้อมูลที่รับส่ง
	opts := options.Find().SetLimit(100).SetProjection(bson.M{"number": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var docs []struct {
		Number string `bson:"number"`
	}
	if err := cursor.All(ctx, &docs); err != nil {
		return
	}

	// 3. หากมีเลขที่ดึงได้ ให้นำไปเพิ่มลงใน Redis Set แบบเรียงลำดับ
	if len(docs) > 0 {
		members := make([]interface{}, len(docs))
		for i, doc := range docs {
			members[i] = doc.Number
		}
		// เพิ่มข้อมูลลงใน Redis Set (SAdd)
		r.redis.SAdd(ctx, redisKey, members...)
		r.redis.Expire(ctx, redisKey, 1*time.Hour) // @TODO: set env (1 ชั่วโมง)
	}
}

func (r *LotteryRepository) UpsertMany(ctx context.Context, tickets []domain.LotteryTicket) error {
	var models []mongo.WriteModel
	now := time.Now()

	for _, t := range tickets {
		doc := fromLotteryDomain(&t)
		if doc.CreatedAt.IsZero() {
			doc.CreatedAt = now
		}
		doc.UpdatedAt = now

		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"number": t.Number}).
			SetUpdate(bson.M{"$set": doc}).
			SetUpsert(true)
		models = append(models, model)
	}

	if len(models) == 0 {
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, models)
	return err
}

func (r *LotteryRepository) MarkAsSold(ctx context.Context, ticketID string, userID string) error {
	filter := bson.M{
		"_id":         ticketID,
		"reserved_by": userID,
		"status":      domain.LotteryStatusReserved,
	}
	update := bson.M{
		"$set": bson.M{
			"status":     domain.LotteryStatusSold,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("ticket not found or not reserved by user")
	}

	return nil
}

func (r *LotteryRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *LotteryRepository) SeedTickets(ctx context.Context, total int) error {
	count, err := r.Count(ctx)
	if err != nil {
		return err
	}

	if count >= int64(total) {
		fmt.Printf("Data already exists (%d tickets). Skipping seed.\n", count)
		return nil
	}

	fmt.Printf("Seeding %d lottery tickets...\n", total)
	startTime := time.Now()

	batchSize := 10000
	for i := 0; i < total; i += batchSize {
		tickets := make([]domain.LotteryTicket, 0, batchSize)
		for j := 0; j < batchSize; j++ {
			num := i + j
			if num >= total {
				break
			}
			ticketNumber := fmt.Sprintf("%06d", num)
			tickets = append(tickets, domain.LotteryTicket{
				Number:    ticketNumber,
				Status:    domain.LotteryStatusAvailable,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		err := r.UpsertMany(ctx, tickets)
		if err != nil {
			return fmt.Errorf("failed to seed batch: %w", err)
		}

		if (i+batchSize)%50000 == 0 || i+batchSize >= total {
			current := i + batchSize
			if current > total {
				current = total
			}
			fmt.Printf("Seeded %d/%d tickets...\n", current, total)
		}
	}

	fmt.Printf("Successfully seeded %d lottery tickets in %v\n", total, time.Since(startTime))
	return nil
}
