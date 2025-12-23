<script lang="ts">
	import { Plus, X } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import { defaultEdge } from "$lib/model.js";
	import type { Edge, Node } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import { type Metadata, type MetadataElement, type MetadataElements, type Parent } from "./api";
	import Header from "./Header.svelte";
	import { type XYPosition, useStore, useSvelteFlow, useUpdateNodeInternals } from "@xyflow/svelte";
	import { History } from "./graph/history.svelte";
	import { resolve } from "$app/paths";
	import DoiInput from "./metadata-panel/DoiInput.svelte";
	import DateInput from "./metadata-panel/DateInput.svelte";
	import OrcidInput from "./metadata-panel/OrcidInput.svelte";
	import Input from "./metadata-panel/Input.svelte";
	import YearInput from "./metadata-panel/YearInput.svelte";
	import AutoComplete from "./metadata-panel/AutoComplete.svelte";
	import { suggestOLS4, suggestSpecies } from "./metadata-panel/auto-complete";
	import type { Snippet } from "svelte";

	let {
		parent,
		metadataElements,
		metadata = $bindable(),
		nodes = $bindable(),
		edges = $bindable(),
		children,
	}: {
		parent?: Parent;
		metadataElements: MetadataElements;
		metadata: Metadata[];
		nodes: Node[];
		edges: Edge[];
		children?: Snippet;
	} = $props();

	const { getViewport } = useSvelteFlow();
	const { width, height } = $derived(useStore());
	const updateNodeInternals = useUpdateNodeInternals();

	const history = new History(nodes, edges, defaultEdge(true), updateNodeInternals);
	$effect(() =>
		// FIXME: https://github.com/sveltejs/svelte/issues/17117
		history.push(
			$state.snapshot(nodes) as unknown as Node[],
			$state.snapshot(edges) as unknown as Edge[],
		),
	);

	let id = nodes.length;
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

	const flatMetadata = $derived.by(() => {
		const map = new Map<string, Omit<MetadataElement, "name">>();
		for (const { name, type, choices, ols4 } of Object.values(metadataElements).flat()) {
			map.set(name, { type, choices, ols4 });
		}
		return map;
	});

	const panelEntries: (Metadata & { mandatory: boolean })[] = $state([]);

	const initialMetadata = $state.snapshot(metadata);
	for (const kind of metadataElements.mandatory) {
		const index = initialMetadata.findIndex((m) => m.kind === kind.name);
		if (index >= 0) {
			const value = initialMetadata[index].value;
			initialMetadata.splice(index, 1);
			panelEntries.push({ mandatory: true, kind: kind.name, value });
		} else {
			panelEntries.push({ mandatory: true, kind: kind.name, value: "" });
		}
	}
	for (const { kind, value } of initialMetadata) {
		panelEntries.push({ mandatory: false, kind, value });
	}

	export function addMetadata() {
		panelEntries.push({ mandatory: false, kind: "", value: "" });
	}

	$effect(() => {
		metadata = panelEntries.map(({ kind, value }) => ({ kind, value }));
	});
</script>

<div class="t-0 absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header {metadataElements}>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="add node"
			onclick={() => addNode()}
		>
			<Plus />
			<span class="pl-2">node</span>
		</button>
		{@render children?.()}
	</Header>
</div>

<div class="h-dvh w-dvw">
	<Graph bind:nodes bind:edges {addNode} onmoveend={() => (offset = 0)}>
		{#snippet panel()}
			<PanelComponent
				onaddmetadata={() => addMetadata()}
				history={{
					undo: history.canUndo() ? () => ([nodes, edges] = history.undo()) : undefined,
					redo: history.canRedo() ? () => ([nodes, edges] = history.redo()) : undefined,
				}}
			>
				<div class="table">
					<div class="table-row">
						<h2 class=" table-cell px-2 text-xl font-bold">Metadata</h2>
						{#if parent}
							<span
								class="table-cell max-w-lg overflow-hidden px-2 text-ellipsis whitespace-nowrap select-text"
							>
								based on
								<a
									class=" text-blue-700 outline-indigo-400 hover:underline focus-visible:outline-2 dark:text-blue-400"
									title={parent.title}
									href={resolve(`/view/${parent.uuid}`)}
								>
									"{parent.title}"
								</a>
							</span>
						{/if}
					</div>
					<div class="table-row">
						<span class="table-cell px-2 py-2 font-bold">kind</span>
						<span class="table-cell px-2 py-2 font-bold">value</span>
					</div>
					{#each panelEntries as m, idx (idx)}
						{@const { type, choices, ols4 } = flatMetadata.get(m.kind) || {}}

						<div class="table-row">
							<div class="table-cell align-middle">
								{#if m.mandatory}
									<div class="px-2 py-1">{m.kind}</div>
								{:else}
									<select
										bind:value={m.kind}
										title={m.kind}
										class="field-sizing-content h-8 rounded px-2 py-1 outline-indigo-400 focus-visible:outline-2"
									>
										<option class="hidden" value="" disabled selected>-- select a kind --</option>
										{#each Object.entries(metadataElements) as [group, values] (group)}
											<optgroup class="dark:bg-neutral-700" label={group}>
												{#each values as { name: kind } (kind)}
													<option class="dark:bg-neutral-700" value={kind}>{kind}</option>
												{/each}
											</optgroup>
										{/each}
									</select>
								{/if}
							</div>

							<div class="table-cell py-0.75 pl-2">
								{#if m.kind === ""}
									<Input value="" disabled placeholder="-- select a kind --" />
								{:else if type === "choice"}
									<select
										bind:value={m.value}
										title={m.value}
										class="field-sizing-content h-8 w-full cursor-pointer rounded border border-neutral-600 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
										required
									>
										<option class="hidden" value="" disabled selected>-- select an option --</option
										>
										{#each choices! as choice (choice)}
											<option class="cursor-pointer dark:bg-neutral-700" value={choice}
												>{choice}</option
											>
										{/each}
									</select>
								{:else if type === "doi"}
									<DoiInput bind:value={m.value} />
								{:else if type === "date"}
									<DateInput bind:value={m.value} />
								{:else if type === "orcid"}
									<OrcidInput bind:value={m.value} />
								{:else if type === "url"}
									<Input placeholder={type} bind:value={m.value} type="url" />
								{:else if type === "year"}
									<YearInput {type} bind:value={m.value} />
								{:else if type === "species"}
									<AutoComplete placeholder={type} suggest={suggestSpecies} bind:value={m.value} />
								{:else if type === "ols4"}
									<AutoComplete
										placeholder={`subclass of '${ols4!.label}' in ontology '${ols4!.onto}'`}
										suggest={(value) => suggestOLS4(value, ols4!)}
										bind:value={m.value}
									/>
								{:else}
									<textarea
										class="field-sizing-content min-h-8 w-full max-w-lg rounded border border-neutral-600 py-1 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
										rows="1"
										bind:value={m.value}
									></textarea>
								{/if}
							</div>

							{#if !m.mandatory}
								<div class="relative table-cell">
									<button
										class="absolute top-2 left-2 inline h-6 cursor-pointer rounded-4xl outline-red-500 hover:bg-red-500/25 focus-visible:outline-2"
										title="delete this metadata entry"
										onclick={() => panelEntries.splice(idx, 1)}
									>
										<X class="stroke-red-500" />
									</button>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>

<svelte:window {onkeydown} />
