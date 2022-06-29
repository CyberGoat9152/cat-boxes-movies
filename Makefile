deploy:
	cd api;\
	$(MAKE) install;\
	$(MAKE) build;\
	../.build/cat-boxes-movies