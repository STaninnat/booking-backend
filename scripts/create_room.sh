#!/bin/bash

curl -X POST "http://localhost:8080/v1/rooms" \
     -H "Content-Type: application/json" \
     -d '{
           "room_name": "Deluxe Room",
           "description": "Lorem ipsum dolor sit amet consectetur adipisicing elit. Distinctio quis eos quod at animi excepturi officia, sed voluptate iure quaerat rem dolorem, nemo corrupti libero ad tempore eius nihil delectus?",
           "price": 10000,
           "max_guests": 6
         }'

echo ""

curl -X POST "http://localhost:8080/v1/rooms" \
     -H "Content-Type: application/json" \
     -d '{
           "room_name": "Middle Room",
           "price": 1002,
           "max_guests": 3
         }'

echo ""

curl -X POST "http://localhost:8080/v1/rooms" \
     -H "Content-Type: application/json" \
     -d '{
           "room_name": "Normal Room",
           "price": 500,
           "max_guests": 2
         }'

echo ""
