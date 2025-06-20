<script lang="ts">
	import { Plus, Save } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import { saveModel } from "$lib/model.js";
	import type { Edge, Node } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import { goto } from "$app/navigation";
	import type { Metadata, MetadataElements } from "./api";
	import Header from "./Header.svelte";

	let {
		version = undefined,
		identifier = $bindable(),
		identifierEditable,
		metadata = $bindable(),
		metadataElements,
		nodes: initialNodes,
		edges: initialEdges,
	}: {
		version?: string;
		identifier: string;
		identifierEditable: boolean;
		metadata: Metadata[];
		metadataElements: MetadataElements;
		nodes: Node[];
		edges: Edge[];
	} = $props();

	let nodes: Node[] = $state.raw(initialNodes);
	let edges: Edge[] = $state.raw(initialEdges);

	let error = $state<string | undefined>(undefined);
	async function save() {
		try {
			const model = await saveModel(identifier, metadata, nodes, edges, version);
			await goto(`/view/${model.identifier}?version=${model.version}`);
		} catch (e: any) {
			error = e;
		}
	}

	function add() {
		const node: Node = {
			id: nodes.length.toString(),
			type: "custom",
			dragHandle: ".drag-handle",
			position: { x: 100, y: 100 },
			data: {
				body: "new",
				color: ["#ff0000", "#00ff00", "#0000ff"][nodes.length % 3],
				editable: true,
			},
		};
		nodes = [...nodes, node];
	}

	let metadataPanel: ReturnType<typeof MetadataPanel>;
</script>

{#if error}
	<div
		class="t-0 absolute z-20 mt-18 w-full rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700"
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

<div class="t-0 absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header createModel={false}>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 hover:bg-indigo-400 dark:bg-indigo-500"
			title="add new node"
			onclick={add}
		>
			<Plus />
			<span class="pl-1">node</span>
		</button>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 hover:bg-indigo-400 dark:bg-indigo-500"
			title="save model"
			onclick={save}
		>
			<Save />
			<span class="pl-2">save</span>
		</button>
	</Header>
</div>

<div class="h-dvh w-dvw">
	<Graph bind:nodes bind:edges interactive={true}>
		{#snippet panel()}
			<PanelComponent onmetadata={() => metadataPanel.addMetadata()}>
				<MetadataPanel
					bind:this={metadataPanel}
					elements={metadataElements}
					bind:identifier
					{identifierEditable}
					bind:metadata
				/>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
