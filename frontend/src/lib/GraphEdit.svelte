<script lang="ts">
	import { CircleX, Plus, Save } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import { saveModel } from "$lib/model.js";
	import type { Edge, Node } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import { beforeNavigate, goto } from "$app/navigation";
	import type { BackendError, Metadata, MetadataElements } from "./api";
	import Header from "./Header.svelte";
	import Alert from "./Alert.svelte";

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

	let alert: ReturnType<typeof Alert>;
	async function save() {
		try {
			const model = await saveModel(identifier, metadata, nodes, edges, version);
			await goto(`/view/${model.identifier}?version=${model.version}`);
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
		}
	}

	function add() {
		const node: Node = {
			id: nodes.length.toString(),
			type: "custom",
			dragHandle: ".drag-handle",
			position: { x: 100, y: 100 },
			data: {
				body: `node #${nodes.length + 1}`,
				color: ["#E40303", "#FF8C00", "#FFED00", "#008026", "#004CFF", "#732982"][nodes.length % 6],
				editable: true,
			},
		};
		nodes = [...nodes, node];
	}

	beforeNavigate((nav) => {
		if (nav.type === "leave") {
			nav.cancel();
		} else {
			let doNavigate = confirm("Leave site? Unsaved changes will be lost.");
			if (!doNavigate) {
				nav.cancel();
			}
		}
	});

	let metadataPanel: ReturnType<typeof MetadataPanel>;
</script>

<Alert bind:this={alert} />

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
		<a
			href={typeof version === "undefined" ? "/" : `/view/${identifier}?version=${version}`}
			class="flex h-10 cursor-pointer items-center rounded-md bg-red-300 px-3 hover:bg-red-400 dark:bg-red-500"
		>
			<CircleX />
			<span class="pl-2">cancel</span>
		</a>
	</Header>
</div>

<div class="h-dvh w-dvw">
	<Graph bind:nodes bind:edges interactive={true}>
		{#snippet panel()}
			<PanelComponent onaddmetadata={() => metadataPanel.addMetadata()}>
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
