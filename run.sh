#update environment variables for any changes
chmod +x writeenv.sh
./writeenv.sh

# Create containers and images 
# Containers are connected and configured according to compose.yaml file
docker-compose up -d

# Starting frontend
docker-compose run --rm frontend 

# Starting backend server and initialising database
docker-compose run --rm backend 

#removes frontend-run containers that are created everytime it starts
docker ps -a --filter "name=frontend-run" -q | xargs docker rm
