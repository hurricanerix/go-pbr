lc:
	@/bin/echo "Lines of Code"
	@/bin/echo -n "Go:   "; find . -name '*.go' -exec wc -l {} \; | tr -s ' ' | cut -d ' ' -f 2 | xargs | sed -e 's/\ /+/g' | bc
	@/bin/echo -n "Vert: "; find . -name '*.vert' -exec wc -l {} \; | tr -s ' ' | cut -d ' ' -f 2 | xargs | sed -e 's/\ /+/g' | bc
	@/bin/echo -n "Frag: "; find . -name '*.frag' -exec wc -l {} \; | tr -s ' ' | cut -d ' ' -f 2 | xargs | sed -e 's/\ /+/g' | bc