<script lang="ts">
	import type { Snippet } from "svelte";
	import Input from "./search/Input.svelte";
	import { goto } from "$app/navigation";

	let { value, button }: { value: string; button: Snippet<[{ onclick: () => any }]> } = $props();
	$inspect(value);

	async function submit() {
		let query = new URLSearchParams();
		query.set("identifier", value);
		await goto(`/versions?${query}`);
	}

	async function onkeydown(key: KeyboardEvent) {
		if (key.code === "Enter") {
			await submit();
		}
	}
</script>

<div class="p-4">
	<Input placeholder="search identifier" bind:value {onkeydown} />
</div>

{@render button({ onclick: submit })}
