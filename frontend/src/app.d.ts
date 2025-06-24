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
		entries<K extends PropertyKey, V>(o: Record<K, V>): [K, V][];
		keys<K extends PropertyKey>(obj: Record<K, unknown>): K[];
	}
}

export {};
