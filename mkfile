README: readme.proto summary.sh
	cp readme.proto README
	./summary.sh >> README

clean:
	rm -f README
