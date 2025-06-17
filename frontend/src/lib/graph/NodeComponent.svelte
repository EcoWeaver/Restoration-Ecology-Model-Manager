<script lang="ts">
	import { type NodeProps, Position, useSvelteFlow } from "@xyflow/svelte";
	import type { Node } from "$lib/model";
	import { Handle } from "@xyflow/svelte";
	import { Edit, Palette, Pencil, Trash2 } from "@lucide/svelte";

	const { id, data }: NodeProps<Node> = $props();
	const { updateNodeData, deleteElements } = useSvelteFlow();

	let body = $state(data.body);
	let color = $state(data.color);
	$effect(() => updateNodeData(id, { body, color }));

	function isLight(s: string) {
		const rgb = parseInt(s.slice(1), 16);
		const r = (rgb >> 16) & 0xff;
		const g = (rgb >> 8) & 0xff;
		const b = rgb & 0xff;

		const luma = r * 0.2126 + g * 0.7152 + b * 0.0722;
		return luma > 127;
	}
	const light = $derived(isLight(color));

	function changeColor() {
		const e = document.createElement("input");
		e.type = "color";
		e.value = color;
		e.click();
		e.onchange = () => (color = e.value);
	}

	function deleteSelf() {
		deleteElements({ nodes: [{ id }] });
	}
</script>

<div
	class="relative rounded-lg px-4 py-2"
	class:text-black={light}
	class:text-white={!light}
	style:background-color={color}
>
	{#if data.editable}
		<div class="cursor-text" contenteditable bind:innerText={body}></div>
		<div
			class="absolute top-full left-0 mt-1 h-10 w-max rounded-full bg-neutral-200 p-2 text-black dark:bg-neutral-700 dark:text-white"
		>
			<button class="cursor-pointer" onclick={changeColor} title="Change color"><Palette /></button>
			<button class="cursor-pointer" onclick={deleteSelf} title="Delete"><Trash2 /></button>
		</div>
	{:else}
		<div>{body}</div>
	{/if}
</div>

<Handle type="target" position={Position.Left} />
<Handle type="source" position={Position.Right} />
