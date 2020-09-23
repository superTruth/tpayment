# 环境
1. Golang 1.13.7
2. Docker 19.03.8

# 编译
1. ./deployment/docker-prod/build_binary.sh
2. ./deployment/docker-prod/build_image.sh

# 部署
docker run -d -p 8001:80 bindo123/tpayment:d85d6ea02bd42136dc692eb32cf4504c37641a26  // 根据实际生成的image进行填写


