build:
	./build.py --all
run:
	./build.py --all && ./run.py
dev:
	docker compose down ; ./build.py --all && ./run.py ; docker compose logs --follow controller
