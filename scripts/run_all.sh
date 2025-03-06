#!/bin/bash

# รันไฟล์ .sh ตามลำดับที่ต้องการ
echo "Running buildprod.sh..."
bash ./script/buildprod.sh

echo "Running migrationup.sh..."
bash ./script/migrationup.sh

echo "Running create_room.sh..."
bash ./script/create_room.sh
