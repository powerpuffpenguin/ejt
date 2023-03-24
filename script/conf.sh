Target="ejt"
Docker="king011/envoy-ejt"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v0.0.4"
View=0
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)