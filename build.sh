chmod +x ./gen_ca.sh
./gen_ca.sh
docker build -t vr009/proxy .
docker run -p 8080:8080 vr009/proxy
