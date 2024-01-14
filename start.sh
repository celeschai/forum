# Create containers and images 
# Containers are connected and configured according to compose.yaml file
docker-compose up -d

# Starting frontend
docker-compose run --rm frontend 

# Starting backend server and initialising database
docker-compose run --rm backend 