async function runWorker() {
	while (true) {
		try {
			const res = await fetch("http://localhost:8080/static/self-lock/", {
				cache: "no-store",
			});
			await res.text();
		} catch {}
	}
}

for (let i = 0; i < 10; i++) {
	runWorker();
}
