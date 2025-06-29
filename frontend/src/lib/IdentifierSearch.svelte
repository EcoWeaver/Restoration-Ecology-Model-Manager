<script lang="ts">
	import { goto } from "$app/navigation";
	import Input from "./search/Input.svelte";
	import SearchButton from "./search/Button.svelte";

	let { identifier = $bindable() }: { identifier: string } = $props();

	async function submit() {
		const query = new URLSearchParams();
		query.set("identifier", identifier);
		await goto(`/versions?${query}`);
	}

	async function onkeydown(key: KeyboardEvent) {
		if (key.code === "Enter") {
			await submit();
		}
	}
</script>

<div class="grow focus-within:outline-2 focus-within:outline-indigo-400">
	<Input placeholder="search by identifier" bind:value={identifier} {onkeydown} />
</div>

<SearchButton onclick={submit} />
