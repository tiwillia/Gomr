#Gomr
-----------
###Instructions
1) Modify the configuration yaml file:
```
cp ./gomr.example.yaml ./gomr.yaml
vim gomr.yaml
...
```
2) Build the docker image:
```
docker build $GOPATH/src/github.com/tiwillia/gomr
```

3) Run it
```
docker run <image_id>
```

### Contributing
See the examplePlugin file for an example on adding your own plugin.
