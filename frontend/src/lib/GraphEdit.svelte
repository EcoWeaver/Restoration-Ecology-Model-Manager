<script lang="ts">
	import { CircleX, Plus, Save } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import { defaultEdge, saveModel } from "$lib/model.js";
	import type { Edge, Node } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import { beforeNavigate, goto } from "$app/navigation";
	import type { BackendError, Metadata, MetadataElements } from "./api";
	import Header from "./Header.svelte";
	import Alert from "./Alert.svelte";
	import { clearNavigate, startNavigate } from "./LoadingSpinner.svelte";
	import { type XYPosition, useStore, useSvelteFlow, useUpdateNodeInternals } from "@xyflow/svelte";
	import { History } from "./graph/history.svelte";
	import { resolve } from "$app/paths";

	let {
		uuid,
		metadata = $bindable(),
		metadataElements,
		nodes: initialNodes,
		edges: initialEdges,
	}: {
		uuid?: string;
		metadata: Metadata[];
		metadataElements: MetadataElements;
		nodes: Node[];
		edges: Edge[];
	} = $props();

	const { getViewport } = useSvelteFlow();
	const { width, height } = $derived(useStore());
	const updateNodeInternals = useUpdateNodeInternals();

	let nodes: Node[] = $state.raw(initialNodes);
	let edges: Edge[] = $state.raw(initialEdges);

	const history = new History(initialNodes, initialEdges, defaultEdge(true), updateNodeInternals);
	$effect(() =>
		// FIXME: https://github.com/sveltejs/svelte/issues/17117
		history.push(
			$state.snapshot(nodes) as unknown as Node[],
			$state.snapshot(edges) as unknown as Edge[],
		),
	);

	let alert: ReturnType<typeof Alert>;
	async function save() {
		try {
			startNavigate();
			const res = await saveModel(uuid, metadata, nodes, edges);
			await goto(resolve(`/view/${res.uuid}`));
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
			clearNavigate();
		}
	}

	let id = initialNodes.length;
	let offset = 0;
	function addNode(position?: XYPosition): string {
		const viewport = getViewport();
		const x = (width / 2 - viewport.x) / viewport.zoom + offset;
		const y = (height / 2 - viewport.y) / viewport.zoom + offset;
		offset += 10;

		const thisId = (id++).toString();
		const node: Node = {
			id: thisId,
			type: "custom",
			dragHandle: ".drag-handle",
			position: position ?? { x, y },
			origin: [0.5, 0.5],
			data: {
				body: "new node",
				color: ["#E40303", "#FF8C00", "#FFED00", "#008026", "#004CFF", "#732982"][id % 6],
				editable: true,
			},
		};

		nodes = [...nodes, node];
		return thisId;
	}

	beforeNavigate((nav) => {
		if (nav.type === "leave" || nav.willUnload) {
			nav.cancel();
		} else if (nav.type !== "goto" || nav.to?.route.id !== "/view/[...model]") {
			const doNavigate = confirm("Leave site? Unsaved changes will be lost.");
			if (!doNavigate) {
				nav.cancel();
			}
		}
	});

	function onkeydown(event: KeyboardEvent) {
		if (event.ctrlKey) {
			if (event.key === "z") {
				if (history.canUndo()) {
					event.preventDefault();
					[nodes, edges] = history.undo();
				}
			} else if (event.key === "Z" || event.key === "y") {
				if (history.canRedo()) {
					event.preventDefault();
					[nodes, edges] = history.redo();
				}
			}
		}
	}

	let metadataPanel: ReturnType<typeof MetadataPanel>;
</script>

<div class="t-0 absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header {metadataElements} createModel={false}>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-2 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="add new node"
			onclick={() => addNode()}
		>
			<Plus />
			<span class="px-1">node</span>
		</button>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="save model"
			onclick={save}
		>
			<Save />
			<span class="pl-2">save</span>
		</button>
		<a
			href={resolve(uuid ? `/view/${uuid}` : "/")}
			title="cancel editing"
			class="flex h-10 cursor-pointer items-center rounded-md bg-red-300 px-3 outline-offset-2 outline-red-400 hover:bg-red-400 focus-visible:outline-2 dark:bg-red-500"
			data-sveltekit-preload-data="tap"
			data-sveltekit-preload-code="hover"
		>
			<CircleX />
			<span class="pl-2 select-none">cancel</span>
		</a>
	</Header>
</div>

<Alert bind:this={alert} />

<div class="h-dvh w-dvw">
	<Graph bind:nodes bind:edges {addNode} onmoveend={() => (offset = 0)}>
		{#snippet panel()}
			<PanelComponent
				onaddmetadata={() => metadataPanel.addMetadata()}
				history={{
					undo: history.canUndo() ? () => ([nodes, edges] = history.undo()) : undefined,
					redo: history.canRedo() ? () => ([nodes, edges] = history.redo()) : undefined,
				}}
			>
				<MetadataPanel bind:this={metadataPanel} elements={metadataElements} bind:metadata />
			</PanelComponent>
		{/snippet}
	</Graph>
</div>

<svelte:window {onkeydown} />
