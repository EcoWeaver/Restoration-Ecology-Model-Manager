<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		getBezierPath,
		MarkerType,
		Position,
		useStore,
		useSvelteFlow,
		type EdgeProps,
	} from "@xyflow/svelte";
	import type { Edge } from "$lib/model";
	import { relationElements } from "$lib/api";

	let { id, data, sourceX, sourceY, targetX, targetY, markerEnd }: EdgeProps<Edge> = $props();

	let [path, labelX, labelY] = $derived(
		getBezierPath({
			sourceX,
			sourceY,
			sourcePosition: Position.Right,
			targetX,
			targetY,
			targetPosition: Position.Left,
		}),
	);

	const { updateEdge, getZoom } = useSvelteFlow();

	let value = $state(data!.kind);

	$effect(() => {
		const v = value;
		relationElements.then((re) => {
			updateEdge(
				id,
				{
					markerEnd: re.causal.includes(v) ? { type: MarkerType.ArrowClosed } : undefined,
					data: { kind: value, editable: data!.editable },
				},
				{ replace: true },
			);
		});
	});
</script>

<BaseEdge {id} {path} {markerEnd} />
<EdgeLabel x={labelX} y={labelY}>
	{#if data!.editable}
		<select bind:value>
			{#await relationElements then es}
				{#each Object.entries(es) as [key, values]}
					<optgroup label={key} style:font-size={`${getZoom()}em`}>
						{#each values as value}
							<option {value}>{value}</option>
						{/each}
					</optgroup>
				{/each}
			{/await}
		</select>
	{:else}
		<p>{data!.kind}</p>
	{/if}
</EdgeLabel>
