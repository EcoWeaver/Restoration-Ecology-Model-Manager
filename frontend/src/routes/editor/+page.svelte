<script lang="ts">
	import type { PageProps } from "./$types";

	import { type Node } from "$lib/model/NodeComponent.svelte";
	import { type Edge } from "$lib/model/EdgeComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import {
		putModel,
		type Node as ApiNode,
		type Relation as ApiRelation,
		type Metadata,
		type ModelRequest,
	} from "$lib/api";
	import Model from "$lib/Model.svelte";

	let { data }: PageProps = $props();
	let { metadataElements } = data;

	let id = $state(0);

	let nodes: Node[] = $state.raw([]);
	let edges: Edge[] = $state.raw([]);

	let error = $state<string | undefined>(undefined);
	let identifier = $state("");
	let metadata = $state<Metadata[]>([]);

	function save() {
		error = undefined;

		const outNodes: ApiNode[] = nodes.map((n) => {
			let { body, color } = n.data;
			return {
				body,
				color,
				position: {
					x: Math.round(n.position.x),
					y: Math.round(n.position.y),
				},
			};
		});
		const outRelations: ApiRelation[] = edges
			.filter((e) => e.type == "custom")
			.map((e) => {
				const from = nodes.findIndex((n) => n.id === e.source);
				const to = nodes.findIndex((n) => n.id === e.target)!;

				let { kind } = e.data!;
				return { kind, from, to };
			});

		const req: ModelRequest = {
			identifier,
			metadata,
			model: { nodes: outNodes, relations: outRelations },
		};

		putModel(req).catch((e) => (error = e));
	}

	function add() {
		const i = id++;
		const node: Node = {
			id: i.toString(),
			type: "custom",
			position: { x: 100, y: 100 },
			data: { body: "new", color: ["#ff0000", "#00ff00", "#0000ff"][i % 3], editable: true },
		};
		nodes = [...nodes, node];
	}
</script>

{#if error}
	<div
		class="relative rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700"
		role="alert"
	>
		<span class="block sm:inline">{error}</span>
		<button
			onclick={() => (error = undefined)}
			class="absolute top-0 right-0 bottom-0 cursor-pointer px-5 py-3 font-extrabold"
		>
			X
		</button>
	</div>
{/if}

<div class="h-dvh w-dvw">
	<Model
		bind:nodes
		bind:edges
		interactive={true}
		onconnect={(c) => {
			const edge: Edge = {
				id: `${c.source}-${c.target}`,
				type: "custom",
				source: c.source,
				target: c.target,
				data: {
					kind: "correlated-with",
					editable: true,
				},
			};
			edges = [...edges, edge];
		}}
	>
		{#snippet info()}
			<h1 class="text-2xl">Restoration Ecology Model Manager</h1>
			<button class="hover:bg-neutral-200" onclick={save}>Save</button>
			<button class="hover:bg-neutral-200" onclick={add}>Add</button>
		{/snippet}

		{#snippet panel()}
			<MetadataPanel elements={metadataElements} bind:identifier bind:metadata />
		{/snippet}
	</Model>
</div>
