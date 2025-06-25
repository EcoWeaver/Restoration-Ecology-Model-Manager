<script lang="ts">
	import Header from "$lib/Header.svelte";
	import otto from "$lib/assets/otto.webp";
</script>

<svelte:head>
	<title>restoration ecology model manager</title>
</svelte:head>

<Header />

<div class="img absolute -z-10 flex w-dvw items-center justify-center p-20 pt-16">
	<img class="max-h-full max-w-full" alt="otto the otter" title="otto the otter" src={otto} />
</div>

<style>
	.img {
		height: calc(100vh - var(--spacing) * 18);
	}
</style>
