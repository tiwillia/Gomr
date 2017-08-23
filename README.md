Gomr
-----------

### Docker Instructions
1) Modify the configuration yaml file:
```
cp ./gomr.example.yaml ./gomr.yaml
vim gomr.yaml
...
```
2) Set up a mysql database and add details to configuration

3) Build the docker image:
```
docker build $GOPATH/src/github.com/tiwillia/gomr
```

4) Run it:
```
docker run <image_id>
```

### OpenShift Insructions
1) Modify the configuration yaml file:
```
cp ./gomr.example.yaml ./gomr.yaml
vim gomr.yaml
...

```
2) Process the template and deploy:
```
oc process -f ./openshift/templates/gomr-mysql.json | oc create -f -
```

3) Watch the build happen with:
```
oc logs -f gomr-build-1
```

### Contributing
See the examplePlugin file for an example on adding your own plugin.
