build:
	cp -r assets output/
	marp slide-deck.md --output output/indexes.html
	marp slide-deck.md --output output/indexes.pdf
	marp slide-deck.md --output output/indexes.png

watch:
	marp slide-deck.md -w --output output/indexes.html
