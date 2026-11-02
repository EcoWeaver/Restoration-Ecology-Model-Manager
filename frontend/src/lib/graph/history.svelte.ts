import type { Polar } from "#lib/api.js";
import { defaultEdgeOptions, type Edge, type Node } from "#lib/model.js";
import type { XYPosition } from "@xyflow/svelte";

type Change =
	| { tag: "addOrRemove"; nodes: Map<string, Node>; edges: Map<string, Edge>; remove: boolean }
	| { tag: "changeNodePosition"; id: string; to: XYPosition }
	| { tag: "changeNodeBody"; id: string; to: string }
	| { tag: "changeNodeColor"; id: string; to: string }
	| { tag: "changeEdgePosition"; id: string; to: Polar }
	| { tag: "changeEdgeKind"; id: string; to: string };

export class History {
	nodes: Node[];
	edges: Edge[];
	updateNodeInternals: (nodeId?: string | string[]) => void;

	undoHist: Change[] = $state([]);
	redoHist: Change[] = $state([]);

	constructor(
		nodes: Node[],
		edges: Edge[],
		updateNodeInternals: (nodeId?: string | string[]) => void,
	) {
		this.nodes = nodes;
		this.edges = edges;
		this.updateNodeInternals = updateNodeInternals;
	}

	push(nodes: Node[], edges: Edge[]) {
		const change =
			this.diffAddOrRemove(nodes, edges) ?? this.diffNodeData(nodes) ?? this.diffEdgeData(edges);

		if (change) {
			this.nodes = structuredClone(nodes);
			this.edges = structuredClone(edges);
			if (change !== "collapsed") {
				this.undoHist.push(change);
				this.redoHist.length = 0;
			}
		}
	}

	diffAddOrRemove(nodes: Node[], edges: Edge[]): Change | undefined {
		let remove: boolean | undefined;
		const changedNodes = new Map<string, Node>();
		const changedEdges = new Map<string, Edge>();

		// because only inserts *or* deletions can happen, the length must be unequal for one of them
		// to have happened
		if (this.nodes.length !== nodes.length) {
			const nodesMap = new Map(this.nodes.map((node) => [node.id, node]));

			for (const node of nodes) {
				if (!nodesMap.delete(node.id)) {
					remove = true;
					changedNodes.set(node.id, node);
				}
			}

			if (changedNodes.size === 0) {
				for (const node of nodesMap.values()) {
					remove = false;
					changedNodes.set(node.id, node);
				}
			}
		}

		// because only inserts *or* deletions can happen, the length must be unequal for one of them
		// to have happened
		if (this.edges.length !== edges.length) {
			const edgesMap = new Map(this.edges.map((edge) => [edge.id, edge]));

			for (const edge of edges) {
				if (!edgesMap.delete(edge.id)) {
					remove = true;
					changedEdges.set(edge.id, edge);
				}
			}

			for (const edge of edgesMap.values()) {
				remove = false;
				changedEdges.set(edge.id, edge);
			}
		}

		if (remove !== undefined) {
			return {
				tag: "addOrRemove",
				nodes: structuredClone(changedNodes),
				edges: structuredClone(changedEdges),
				remove,
			};
		}
	}

	diffNodeData(nodes: Node[]): Change | "collapsed" | undefined {
		const nodesMap = new Map(this.nodes.map((node) => [node.id, node]));

		for (const node of nodes) {
			const oldNode = nodesMap.get(node.id)!;
			if (oldNode.position.x !== node.position.x || oldNode.position.y !== node.position.y) {
				// if the last change was also "changeNodePosition", collapse the two events
				if (this.undoHist.length > 0) {
					const undoHistEnd = this.undoHist[this.undoHist.length - 1];
					if (undoHistEnd.tag === "changeNodePosition" && undoHistEnd.id === node.id) {
						return "collapsed";
					}
				}
				return {
					tag: "changeNodePosition",
					id: oldNode.id,
					to: { ...oldNode.position },
				};
			} else if (oldNode.data.body !== node.data.body) {
				// if the last change was also "changeNodeBody", collapse the two events
				if (this.undoHist.length > 0) {
					const undoHistEnd = this.undoHist[this.undoHist.length - 1];
					if (undoHistEnd.tag === "changeNodeBody" && undoHistEnd.id === node.id) {
						return "collapsed";
					}
				}
				return {
					tag: "changeNodeBody",
					id: oldNode.id,
					to: oldNode.data.body,
				};
			} else if (oldNode.data.color !== node.data.color) {
				return {
					tag: "changeNodeColor",
					id: oldNode.id,
					to: oldNode.data.color,
				};
			}
		}
	}

	diffEdgeData(edges: Edge[]): Change | "collapsed" | undefined {
		const edgesMap = new Map(this.edges.map((edge) => [edge.id, edge]));

		for (const edge of edges) {
			const oldEdge = edgesMap.get(edge.id)!;
			const oldData = oldEdge.data ?? defaultEdgeOptions.data!;
			const newData = edge.data ?? defaultEdgeOptions.data!;
			if (
				oldData.position.r !== newData.position.r ||
				oldData.position.phi !== newData.position.phi
			) {
				// if the last change was also "changeEdgePosition", collapse the two events
				if (this.undoHist.length > 0) {
					const undoHistEnd = this.undoHist[this.undoHist.length - 1];
					if (undoHistEnd.tag === "changeEdgePosition" && undoHistEnd.id === edge.id) {
						return "collapsed";
					}
				}
				return {
					tag: "changeEdgePosition",
					id: oldEdge.id,
					to: { ...oldData.position },
				};
			} else if (oldData.kind !== newData.kind) {
				return {
					tag: "changeEdgeKind",
					id: oldEdge.id,
					to: oldData.kind,
				};
			}
		}
	}

	apply(change: Change) {
		switch (change.tag) {
			case "addOrRemove": {
				if (change.remove) {
					this.nodes = this.nodes.filter((node) => !change.nodes.has(node.id));
					this.edges = this.edges.filter((edge) => !change.edges.has(edge.id));
				} else {
					this.nodes = [...this.nodes, ...change.nodes.values().map((e) => structuredClone(e))];
					this.edges = [...this.edges, ...change.edges.values().map((e) => structuredClone(e))];
				}
				for (const { source, target } of change.edges.values()) {
					this.updateNodeInternals([source, target]);
				}
				change.remove = !change.remove;
				break;
			}
			case "changeNodePosition": {
				const node = this.nodes.find((node) => node.id === change.id)!;
				const oldPosition = { ...node.position };
				node.position = change.to;
				change.to = oldPosition;
				break;
			}
			case "changeNodeBody": {
				const node = this.nodes.find((node) => node.id === change.id)!;
				const oldBody = node.data.body;
				node.data.body = change.to;
				change.to = oldBody;
				break;
			}
			case "changeNodeColor": {
				const node = this.nodes.find((node) => node.id === change.id)!;
				const oldColor = node.data.color;
				node.data.color = change.to;
				change.to = oldColor;
				break;
			}
			case "changeEdgePosition": {
				const edge = this.edges.find((edge) => edge.id === change.id)!;
				const oldData = edge.data ?? defaultEdgeOptions.data!;
				const oldPosition = { ...oldData.position };
				edge.data = { ...oldData, position: change.to };
				change.to = oldPosition;
				break;
			}
			case "changeEdgeKind": {
				const edge = this.edges.find((edge) => edge.id === change.id)!;
				const oldData = edge.data ?? defaultEdgeOptions.data!;
				const oldKind = oldData.kind;
				edge.data = { ...oldData, kind: change.to };
				change.to = oldKind;
				break;
			}
		}
	}

	canUndo(): boolean {
		return this.undoHist.length > 0;
	}

	undo(): [Node[], Edge[]] {
		const change = $state.snapshot(this.undoHist.pop()) as Change;
		if (change) {
			this.apply(change);
			this.redoHist.push(change);
		}
		return [structuredClone(this.nodes), structuredClone(this.edges)];
	}

	canRedo(): boolean {
		return this.redoHist.length > 0;
	}

	redo(): [Node[], Edge[]] {
		const change = $state.snapshot(this.redoHist.pop()) as Change;
		if (change) {
			this.apply(change);
			this.undoHist.push(change);
		}
		return [structuredClone(this.nodes), structuredClone(this.edges)];
	}
}
