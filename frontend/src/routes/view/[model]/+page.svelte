<script lang="ts">
	let { data } = $props();
	let { model } = data;

	console.log(model);
</script>

<table>
	<thead>
		<tr>
			<th>kind</th>
			<th>value</th>
		</tr>
	</thead>
	<tbody>
		{#each model.metadata as { kind, value }}
			<tr>
				<td>{kind}</td>
				<td class="max-w-128">{value}</td>
			</tr>
		{/each}
	</tbody>
</table>

<div class="my-2">{JSON.stringify(model.model)}</div>
