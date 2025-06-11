export function groupByDateBuckets<T>(
	items: T[],
	getDate: (item: T) => Date
): Record<'Today' | 'Yesterday' | 'Last 7 Days' | 'Older', T[]> {
	const norm = (d: Date) => new Date(d.getFullYear(), d.getMonth(), d.getDate());
	const today = norm(new Date());
	const yesterday = new Date(today);
	yesterday.setDate(today.getDate() - 1);
	const last7 = new Date(today);
	last7.setDate(today.getDate() - 7);

	const buckets = {
		Today: [] as T[],
		Yesterday: [] as T[],
		'Last 7 Days': [] as T[],
		Older: [] as T[],
	};

	items.forEach(item => {
		const day = norm(getDate(item));
		if (day >= today) buckets.Today.push(item);
		else if (day >= yesterday) buckets.Yesterday.push(item);
		else if (day >= last7) buckets['Last 7 Days'].push(item);
		else buckets.Older.push(item);
	});

	return buckets;
}

export function formatDateAsString(d: Date | string): string {
	return new Date(d).toLocaleDateString('en-US', {
		day: '2-digit',
		month: 'short',
		year: 'numeric',
	});
}
