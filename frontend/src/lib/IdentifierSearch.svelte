<script lang="ts">
	import { goto } from "$app/navigation";
	import Input from "./search/Input.svelte";
	import SearchButton from "./search/Button.svelte";

	let { identifier }: { identifier: string } = $props();

	async function submit() {
		let query = new URLSearchParams();
		query.set("identifier", identifier);
		await goto(`/versions?${query}`);
	}

	async function onkeydown(key: KeyboardEvent) {
		if (key.code === "Enter") {
			await submit();
		}
	}
</script>

<div class="grow px-0 py-4">
	<Input placeholder="search identifier" bind:value={identifier} {onkeydown} />
</div>

<SearchButton onclick={submit} />
