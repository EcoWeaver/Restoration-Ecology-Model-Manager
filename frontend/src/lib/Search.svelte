<script lang="ts">
	import { goto } from "$app/navigation";
	import { metadataElements } from "./api";
	import Thing from "./search/Thing.svelte";

	let input = $state("");

	let elements = metadataElements.then((e) =>
		Object.values(e)
			.flat()
			.map((e) => e.name),
	);
	let suggestions = $state<string[]>([]);

	let focused = $state(false);
	let thing = $state<string | null>(null);
	async function key(key: KeyboardEvent) {
		if (key.code === "Enter") {
			if (thing === null && suggestions[0]) {
				thing = suggestions[0];
				input = "";
			} else {
				await goto(`/search?${thing}=${input}`);
				thing = "";
			}
		}
	}

	$effect(() => {
		if (thing) {
			suggestions = [];
			return;
		}

		input;
		elements.then((elements) => {
			if (input === "") {
				suggestions = elements;
			} else {
				suggestions = elements.filter((elem) => elem.startsWith(input));
			}
		});
	});
</script>

<div class="p-4">
	<span
		class="flex h-10 w-96 rounded outline-1 focus-within:outline-2 focus-within:outline-blue-500"
	>
		{#if thing}
			<Thing {thing} />
		{/if}

		<input
			class="block grow px-2 py-1 focus-within:outline-none"
			onkeydown={key}
			bind:focused
			bind:value={input}
		/>
	</span>

	{#if focused}
		{#each suggestions as suggestion}
			<p>{suggestion}</p>
		{/each}
	{/if}
</div>
