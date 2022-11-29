source scripts/values.sh

echo -e "\nBuilding as version ${c}${b}$VERSION${n}"

docker build -t $IMAGE_NAME:$VERSION .

echo -e "\n${b}Compressing image...${n}"
docker save $IMAGE_NAME:$VERSION | gzip > $BUILD_DIR/${IMAGE_NAME}_${VERSION}.tar.gz
echo -e "Saved to ${b}${g}${BUILD_DIR}/${IMAGE_NAME}_${VERSION}.tar.gz${n}"
