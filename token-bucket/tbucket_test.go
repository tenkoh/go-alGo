package tbucket_test

// 1. bucketに空きがある場合はtokenを取り出せる
// 2. bucketに空きがない場合はtokenを取り出せない(ブロックされる)
// 3. bucketからの取り出しはキャンセルできる
// 4. bucketには決めた間隔でtokenが補充される
// 5. bucketには最大値を超えてtokenを補充できない
