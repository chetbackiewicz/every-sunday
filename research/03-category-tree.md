# Hierarchical Category Tree UI Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: react-arborist ✅ SELECTED

---

## Recommendation Summary

**react-arborist** chosen because:
1. **Headless architecture** - Matches TanStack Table approach for consistent codebase
2. **Built-in drag & drop** - No need to integrate separate DnD library for reordering
3. **Easy depth limiting** - `node.level` property makes 2-level limit trivial
4. **Free and open source** - No licensing costs
5. **Performance** - Virtual scrolling built-in for large category lists
6. **Modern React patterns** - Hooks-based API aligns with project stack
7. **Sufficient features** - Everything needed without bloat

**Alternative if react-arborist doesn't work**: react-dnd-treeview (simpler API, good docs)

**Avoid**:
- **MUI X Tree View** - Paid license required for drag & drop, opinionated styling
- **react-beautiful-dnd** - Archived, no longer maintained
- **Custom solution** - Unnecessary complexity when good libraries exist

---

## react-arborist - Modern Headless Tree Component

**Source**: https://github.com/brimdata/react-arborist

### Key Features

- ✅ **Headless architecture** - Full control over styling and rendering
- ✅ **Drag & drop built-in** - Using react-dnd with HTML5 backend
- ✅ **Node depth tracking** - `node.level` property for depth enforcement
- ✅ **Expand/collapse** - `tree.open()`, `tree.close()`, `tree.toggle()` APIs
- ✅ **Virtual scrolling** - Uses react-window for performance with large trees
- ✅ **Create/Delete nodes** - `onCreate`, `onDelete` handlers
- ✅ **TypeScript first** - Excellent type support
- ✅ **Custom node rendering** - Render prop pattern with `NodeRendererProps`
- ✅ **Selection support** - Single and multi-select built-in
- ✅ **Search/filter** - `searchTerm` and `searchMatch` props
- ✅ **Controlled state** - Can control `openState` externally
- ✅ **Keyboard navigation** - Arrow keys, Enter, Space

### Depth Limiting Implementation

```typescript
<Tree
  data={categories}
  disableDrop={(args) => {
    const { parentNode, dragNodes } = args
    // Block drops if parent depth >= 2 (allowing Parent→Child→Grandchild)
    return parentNode.level >= 2
  }}
  onCreate={({ parentId, parentNode, index, type }) => {
    // Block creation if parent depth >= 2
    if (parentNode && parentNode.level >= 2) {
      return null // Prevent node creation
    }
    return { id: newId, name: '' }
  }}
/>
```

### Node Rendering with + Buttons

```typescript
function CategoryNode({ node, style, tree }: NodeRendererProps<Category>) {
  const canAddChild = node.level < 2 // Only show + if not at max depth
  
  return (
    <div style={style}>
      {node.isInternal && (
        <button onClick={() => node.toggle()}>
          {node.isOpen ? '▼' : '▶'}
        </button>
      )}
      
      <span>{node.data.name}</span>
      
      {canAddChild && (
        <button onClick={() => tree.onCreate({
          parentId: node.id,
          parentNode: node,
          index: node.children?.length || 0,
          type: 'internal'
        })}>
          + Add
        </button>
      )}
    </div>
  )
}
```

### Pros

- Modern, actively maintained (2023+)
- Perfect depth control via `node.level`
- Built-in drag & drop with depth restrictions
- Headless - matches TanStack Table philosophy
- Virtual scrolling for performance
- Simple API for CRUD operations

### Cons

- Smaller community than alternatives
- Requires react-dnd (additional dependency)
- Less comprehensive docs than MUI

---

## Alternative Options (Not Selected)

### react-dnd-treeview

**Source**: https://github.com/minop1205/react-dnd-treeview

**Pros**:
- Simple flat data structure
- Excellent drag & drop between multiple trees
- Good examples and documentation
- Actively maintained
- Flexible node rendering

**Cons**:
- Manual depth calculation required
- No built-in virtual scrolling
- Must manage `droppable` flag manually for depth limiting

### MUI X Tree View

**Source**: https://mui.com/x/react-tree-view/

**Why Not Selected**:
- **Paid license required for key features** - Drag & drop, lazy loading, virtualization all in Pro ($200+/dev/year)
- **Opinionated Material Design styling** - Harder to customize for custom budget UI
- **Large bundle size** - Full MUI ecosystem required
- **May be overkill** - Many features not needed for this use case

### react-beautiful-dnd

**Source**: https://github.com/atlassian/react-beautiful-dnd

**Status**: ⚠️ **ARCHIVED** - No longer maintained as of August 2025
- Atlassian recommends migrating to Pragmatic drag and drop
- Still popular but no future updates or bug fixes
- **Do not use for new projects**

---

## Comparison Matrix

| Feature            | react-arborist    | react-dnd-treeview | MUI X Tree View               | Custom Solution |
| ------------------ | ----------------- | ------------------ | ----------------------------- | --------------- |
| License            | MIT (Free)        | MIT (Free)         | Community (Free) / Pro (Paid) | N/A             |
| Bundle Size        | ~20kb             | ~15kb              | ~100kb+ (with MUI)            | Minimal         |
| Drag & Drop        | Built-in ✅        | Built-in ✅         | Pro only 💰                    | Custom          |
| Depth Control      | Easy (node.level) | Manual             | Easy                          | Easy            |
| Styling Control    | 100% Custom       | 100% Custom        | Material Design + Overrides   | 100% Custom     |
| Virtual Scrolling  | Built-in ✅        | ❌                  | Pro only 💰                    | Custom          |
| TypeScript         | Excellent         | Good               | Excellent                     | N/A             |
| Learning Curve     | Medium            | Low                | Medium-High                   | High            |
| Active Maintenance | ✅ Active          | ✅ Active           | ✅ Active                      | N/A             |
