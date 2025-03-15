package main

// "github.com/go-redis/redis/v8"

// var (
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379", // use your Redis server address
// 	})
// 	val, err := redisClient.Get(ctx, cacheKey).Result()
// )

// func GetStorageLimit(userID int64) (int64, error) {
// 	cacheKey := fmt.Sprintf("user:%d:storage_limit", userID)

// 	// Try Redis first
// 	val, err := redisClient.Get(cacheKey).Result()
// 	if err == nil {
// 		limit, _ := strconv.ParseInt(val, 10, 64)
// 		return limit, nil
// 	redisClient.Set(ctx, cacheKey, storageLimit, 10*time.Minute)

// 	// If not in cache, query the DB
// 	var storageLimit int64
// 	err = db.QueryRow("SELECT storage_limit FROM users WHERE id = $1", userID).Scan(&storageLimit)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Store in Redis with an expiry
// 	redisClient.Set(cacheKey, storageLimit, 10*time.Minute)

// 	return storageLimit, nil
// }
