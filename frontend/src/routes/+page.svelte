<script lang="ts">
	import Header from "#lib/Header.svelte";
	import otto from "#lib/assets/otto.webp";
</script>

<svelte:head>
	<title>restoration ecology model manager</title>
</svelte:head>

<div class="absolute w-full">
	<Header newModel={true} />
</div>

<div class="absolute -z-10 flex h-full w-full justify-center p-20 pt-34">
	<img class="max-h-full max-w-full" alt="otto the otter" title="otto the otter" src={otto} />
</div>
