export function stripURL(value: string, url: string): string {
	value = value.trim();
	if (value.startsWith(url)) {
		value = value.slice(url.length);

		if (value.endsWith("/")) {
			value = value.slice(0, -1);
		}
	}

	return value;
}
