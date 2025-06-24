// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}

	interface ObjectConstructor {
		entries<TKey extends PropertyKey, TVal>(o: Record<TKey, TVal>): [TKey, TVal][];
		keys<TKey extends PropertyKey>(obj: Record<TKey, unknown>): TKey[];
	}
}

export {};
