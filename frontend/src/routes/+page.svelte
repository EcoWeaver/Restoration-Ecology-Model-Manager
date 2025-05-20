<script lang="ts">
	import { Background, BackgroundVariant, Panel, SvelteFlow } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import NodeComponent, { type Node } from '$lib/NodeComponent.svelte';
	import EdgeComponent, { type Edge } from '$lib/EdgeComponent.svelte';
	import MetadataPanel from '$lib/MetadataPanel.svelte';
	import {
		metadataElements,
		type Node as ApiNode,
		type Relation as ApiRelation,
		type Graph as ApiGraph,
		type Metadata,
		type GraphRequest
	} from '$lib/api';

	let nodes: Node[] = $state.raw([
		{
			id: '1',
			type: 'custom',
			position: { x: 100, y: 100 },
			data: { body: 'plant', color: '#ff0000' }
		},
		{
			id: '2',
			type: 'custom',
			position: { x: 300, y: 200 },
			data: { body: 'tree', color: '#00ff00' }
		},
		{
			id: '3',
			type: 'custom',
			position: { x: 100, y: 300 },
			data: { body: 'animal', color: '#0000ff' }
		}
	]);

	let edges: Edge[] = $state.raw([
		{ id: '1', type: 'custom', source: '1', target: '2', data: { kind: 'correlated-with' } }
	]);

	let metadata = $state<Metadata[]>([]);

	function save() {
		const outNodes: ApiNode[] = nodes.map((n) => ({ ...n.data, position: n.position }));
		const outRelations: ApiRelation[] = edges.map((e) => {
			const from = nodes.findIndex((n) => n.id === e.source);
			const to = nodes.findIndex((n) => n.id === e.target)!!;
			return {
				...e.data!!,
				from,
				to
			};
		});

		const req: GraphRequest = {
			metadata,
			graph: { nodes: outNodes, relations: outRelations }
		};
		console.log(JSON.stringify(req, null, 2));
	}
</script>

<div class="h-dvh w-dvw">
	<SvelteFlow
		bind:nodes
		bind:edges
		proOptions={{ hideAttribution: true }}
		nodeTypes={{ custom: NodeComponent }}
		edgeTypes={{ custom: EdgeComponent }}
	>
		<Panel position="top-left">
			<h1 class="text-2xl">Restoration Ecology Model Manager</h1>
			<button onclick={save}>Save</button>
		</Panel>
		<Panel position="top-right">
			{#await metadataElements then elements}
				<MetadataPanel {elements} bind:metadata />
			{/await}
		</Panel>
		<Background variant={BackgroundVariant.Dots} />
	</SvelteFlow>
</div>
