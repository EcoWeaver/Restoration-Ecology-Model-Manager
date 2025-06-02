<script lang="ts">
	import Search from "$lib/Search.svelte";
</script>

<div class="flex">
	<Search />
	<a
		href="./editor"
		class="bg m-4 ml-0 flex h-10 place-items-center rounded bg-indigo-300 px-4 hover:bg-indigo-400"
	>
		<span class="font-medium">+ new model</span>
	</a>
</div>
