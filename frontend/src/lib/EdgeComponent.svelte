<script lang="ts" module>
	export type Edge = FlowEdge<Omit<Relation, 'from' | 'to'>, 'custom'>;
</script>

<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		getBezierPath,
		MarkerType,
		type Edge as FlowEdge,
		Position,
		useSvelteFlow,
		type EdgeProps
	} from '@xyflow/svelte';
	import { relationElements, type Relation } from '$lib/api';

	let { id, data, sourceX, sourceY, targetX, targetY, markerEnd }: EdgeProps<Edge> = $props();

	let [path, labelX, labelY] = $derived(
		getBezierPath({
			sourceX,
			sourceY,
			sourcePosition: Position.Right,
			targetX,
			targetY,
			targetPosition: Position.Left
		})
	);

	const { updateEdge } = useSvelteFlow();

	let value = $state(data.kind);

	$effect(() => {
		const v = value;
		relationElements.then((re) => {
			console.log(re.causal.includes(v));
			updateEdge(
				id,
				{
					markerEnd: re.causal.includes(v) ? { type: MarkerType.ArrowClosed } : undefined,
					data: { relation: value }
				},
				{ replace: true }
			);
		});
	});
</script>

<BaseEdge {id} {path} {markerEnd} />
<EdgeLabel x={labelX} y={labelY}>
	<select bind:value>
		{#await relationElements then es}
			{#each Object.entries(es) as [key, values]}
				<optgroup label={key}>
					{#each values as value}
						<option {value}>{value}</option>
					{/each}
				</optgroup>
			{/each}
		{/await}
	</select>
</EdgeLabel>
