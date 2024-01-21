#update environment variables for any changes
chmod +x writeenv.sh
./writeenv.sh

# Create containers and images 
# Containers are connected and configured according to compose.yaml file
docker-compose up -d

#removes frontend-run containers that are created everytime it starts
docker ps -a --filter "name=frontend-run" -q | xargs docker rm
