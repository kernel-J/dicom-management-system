# dicom-management-system

## Requirements

1. Docker https://docs.docker.com/get-docker/
2. Docker Compose https://docs.docker.com/compose/install/

## Routes

ROUTE | METHOD | FUNCTION | PARAMETERS 
--- | --- | --- | --- 
/upload | POST | Saves a file to the local storage |  file: path-to-file 
/dicom/{id}/{tag} | GET | Extracts and returns attribute of a DICOM image | id: id of the file, tag: tag
/convert/{id} | GET | Converts a dicom file into png and returns the file | id: id of the file

## Starting the application

To start the application run the following command from the root directory

```Bash
docker compose up dicom-management-service
```

## Usage
To upload a DICOM file
_This gives back a `{file id}.dcm` in the a response_
```Bash
curl -X POST -F "file=@path-to-file" http://localhost:8080/upload 
```

To extract attributes of a DICOM file
`tag` should be in the format format (XXXX,XXXX)
```Bash
curl http://localhost:8080/dicom/{file id}/{tag}

# Example
curl http://localhost:8080/dicom/88053a46-5b4a-4ad6-ae26-2e1d19f5c413/\(0040,0275\)
```

To convert the DICOM image to a PNG
```Bash
curl http://localhost:8080/convert/{file id} --output filename

# Example
curl http://localhost:8080/convert/88053a46-5b4a-4ad6-ae26-2e1d19f5c413 --output converted.png
```

## Testing
To run unit tests
```
docker-compose up test
```
