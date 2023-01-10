build:
	./build.py --all --parallel
run:
	./build.py --all --parallel && ./run.py
dev:
	docker-compose down ; ./build.py --all --parallel && ./run.py ; docker-compose logs --follow controller
