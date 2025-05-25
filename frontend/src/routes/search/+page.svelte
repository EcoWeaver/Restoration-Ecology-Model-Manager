<script lang="ts">
	let { data } = $props();
	let { item } = data;

	console.log(item);
</script>

<div>{JSON.stringify(item)}</div>
