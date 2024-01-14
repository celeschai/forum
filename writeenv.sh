if [ -f .env ]; then
    cp .env backend/.env
        echo "Copied .env file to backend" 
    
    envFilePath="frontend/.env"
    cp .env "$envFilePath"
    
    # Prefix all variable names with REACT_APP_
    awk -F= '!/^#/ && NF==2 { if ($1 == "HTTPS") print $1 "=" $2; else print "REACT_APP_" $1 "=" $2 }' "$envFilePath" > "$envFilePath.tmp"
    mv "$envFilePath.tmp" "$envFilePath"
        echo "Copied modified .env file to frontend" 
else
    echo ".env file not found"
fi

