# Multi-Month Management - Month Switching UX Research

**Status**: ✅ Complete  
**Date**: January 19, 2025  
**Decision**: Controlled `<select>` with useState + localStorage persistence

---

## Recommendation Summary

**Chosen Approach: Controlled Select Dropdown**

Based on React documentation and UX patterns, the month selector should use a simple controlled `<select>` component with `useState` for immediate UI updates and `useEffect` for localStorage persistence.

**Why This Approach**:
- ✅ **Simple and familiar** - Users understand dropdowns for month selection
- ✅ **Accessible** - Native `<select>` has built-in keyboard navigation
- ✅ **No routing complexity** - State-based, doesn't require React Router
- ✅ **localStorage integration** - Easy persistence with useEffect pattern
- ✅ **Optimal for desktop** - Budget tracking is primarily desktop workflow
- ✅ **Future-proof** - Can add prev/next buttons alongside dropdown

**Alternative Considered**: URL-based routing (`/budget/2025-10`)
- Better for sharing/bookmarking specific months
- Requires React Router dependency
- More complex initial implementation
- **Recommendation**: Add later if backend API is built

---

## Core Requirements

From spike document:
- "The user should have different months in a drop down that they can create a budget tracking sheet"
- Must save current month data before switching
- Each month has independent budget categories and CSV files
- Need to load saved data when switching back to a month

---

## UI Pattern Options

### Option 1: Select Dropdown (RECOMMENDED)

**Implementation Pattern**:
```typescript
function MonthSelector({ 
  currentMonth, 
  availableMonths, 
  onMonthChange 
}: MonthSelectorProps) {
  return (
    <label>
      Budget Month:
      <select 
        value={currentMonth} 
        onChange={(e) => onMonthChange(e.target.value)}
      >
        {availableMonths.map(month => (
          <option key={month} value={month}>
            {formatMonth(month)} {/* e.g., "October 2025" */}
          </option>
        ))}
      </select>
    </label>
  )
}

// Helper function to format month display
function formatMonth(monthString: string): string {
  const [year, month] = monthString.split('-')
  const date = new Date(parseInt(year), parseInt(month) - 1)
  return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}
```

**Pros**:
- Native HTML element - no dependencies
- Accessible with keyboard (arrow keys, type-ahead)
- Familiar UX pattern
- Small visual footprint
- Works on all devices

**Cons**:
- Less visually prominent than tabs
- No calendar view of available months

---

### Option 2: Tab Navigation

**Example Libraries**: 
- Headless UI Tabs
- React Aria Tabs
- Custom implementation

**Pros**:
- More prominent visual presence
- Easy to see all available months at glance
- Click anywhere on tab to switch

**Cons**:
- Takes more horizontal space
- Limited to ~6-8 visible tabs before scrolling needed
- Requires additional component library OR custom implementation
- Less intuitive for many months

**When to Use**: If user typically works with only 2-4 months at a time

---

### Option 3: Calendar Picker

**Example Libraries**:
- react-datepicker (month-only mode)
- @mui/x-date-pickers
- react-day-picker

**Pros**:
- Visual calendar interface
- Easy to jump to any month/year
- Intuitive for date-based selection

**Cons**:
- Requires date picker library (~15-50kb)
- More clicks to select month (open picker → select month → select year → confirm)
- Overkill for simple month selection

**When to Use**: If app expands to support many years of historical data

---

### Option 4: Prev/Next Buttons with Display

**Implementation Pattern**:
```typescript
function MonthNavigation({ currentMonth, onMonthChange }: MonthNavProps) {
  const handlePrev = () => {
    const prev = subtractMonth(currentMonth)
    onMonthChange(prev)
  }
  
  const handleNext = () => {
    const next = addMonth(currentMonth)
    onMonthChange(next)
  }
  
  return (
    <div className="month-nav">
      <button onClick={handlePrev} aria-label="Previous month">
        ←
      </button>
      <span className="current-month">
        {formatMonth(currentMonth)}
      </span>
      <button onClick={handleNext} aria-label="Next month">
        →
      </button>
    </div>
  )
}
```

**Pros**:
- Quick navigation between adjacent months
- Minimal UI
- Clear current month display

**Cons**:
- Slow for jumping multiple months
- No overview of available months
- Assumes sequential month access

**Recommendation**: **Combine with dropdown** - Best of both worlds
```typescript
function MonthControls() {
  return (
    <div className="month-controls">
      <button onClick={handlePrev}>←</button>
      <select value={currentMonth} onChange={handleMonthChange}>
        {/* month options */}
      </select>
      <button onClick={handleNext}>→</button>
    </div>
  )
}
```

---

## State Management Pattern

### Current Month State

```typescript
interface AppState {
  currentMonth: string  // Format: "YYYY-MM" (e.g., "2025-10")
  monthData: Map<string, MonthBudgetData>  // Keyed by month string
}

interface MonthBudgetData {
  categories: CategoryNode[]
  csvFiles: CSVFile[]
  // ... other month-specific data
}

// Component state
function BudgetApp() {
  const [currentMonth, setCurrentMonth] = useState<string>(() => {
    // Load from localStorage on initial render
    return localStorage.getItem('currentMonth') || getCurrentMonth()
  })
  
  const [monthData, setMonthData] = useState<Map<string, MonthBudgetData>>(() => {
    // Load all saved months from localStorage
    const saved = localStorage.getItem('monthData')
    return saved ? new Map(JSON.parse(saved)) : new Map()
  })
  
  // Persist current month to localStorage when it changes
  useEffect(() => {
    localStorage.setItem('currentMonth', currentMonth)
  }, [currentMonth])
  
  // Persist month data whenever it changes
  useEffect(() => {
    localStorage.setItem('monthData', JSON.stringify([...monthData]))
  }, [monthData])
}
```

**Key Insight from React docs**: 
- Use `useState` lazy initializer (function) to read localStorage only once on mount
- Use `useEffect` to synchronize state changes to localStorage
- Don't read localStorage during every render (performance issue)

---

## Month Switching Behavior

### Before Switching: Auto-save Current Month

```typescript
function handleMonthChange(newMonth: string) {
  // Save current month data before switching
  setMonthData(prev => {
    const updated = new Map(prev)
    updated.set(currentMonth, {
      categories: currentCategories,
      csvFiles: currentCsvFiles
      // ... capture all current state
    })
    return updated
  })
  
  // Switch to new month
  setCurrentMonth(newMonth)
  
  // Load new month data (or create empty if first time)
  const newMonthData = monthData.get(newMonth) || createEmptyMonthData()
  // ... restore state for new month
}
```

### After Switching: Load or Initialize New Month

```typescript
useEffect(() => {
  // When currentMonth changes, load that month's data
  const data = monthData.get(currentMonth)
  
  if (data) {
    // Restore saved data
    setCategories(data.categories)
    setCsvFiles(data.csvFiles)
  } else {
    // First time visiting this month - initialize empty budget
    setCategories(createDefaultCategories())
    setCsvFiles([])
  }
}, [currentMonth])
```

---

## URL Routing Alternative (Future Enhancement)

If backend API is built, consider adding URL-based routing:

```typescript
// Using React Router v6
import { useParams, useNavigate } from 'react-router-dom'

function BudgetPage() {
  const { month } = useParams<{ month: string }>()  // e.g., "2025-10"
  const navigate = useNavigate()
  
  const handleMonthChange = (newMonth: string) => {
    navigate(`/budget/${newMonth}`)
  }
  
  // ... rest of component
}

// Route definition
<Route path="/budget/:month" element={<BudgetPage />} />
```

**Benefits**:
- Shareable URLs: `example.com/budget/2025-10`
- Browser back/forward buttons work
- Bookmark specific months

**Tradeoffs**:
- Requires React Router dependency
- More complex state management (sync URL with state)
- Must handle invalid month params

**Recommendation**: Add when backend API exists, keep localStorage as fallback

---

## Available Months List

### Dynamic Month Generation

```typescript
function generateAvailableMonths(startDate: Date, monthCount: number = 12): string[] {
  const months: string[] = []
  
  for (let i = 0; i < monthCount; i++) {
    const date = new Date(startDate)
    date.setMonth(date.getMonth() - i)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    months.push(`${year}-${month}`)
  }
  
  return months
}

// Usage
const availableMonths = generateAvailableMonths(new Date(), 12)  
// Returns: ["2025-10", "2025-09", "2025-08", ...]
```

### OR: Track User-Created Months

```typescript
interface AppState {
  createdMonths: Set<string>  // Only months user has worked on
}

function handleCreateMonth(monthString: string) {
  setCreatedMonths(prev => new Set(prev).add(monthString))
}
```

**Recommendation**: **Generate fixed range** (e.g., current month + 11 previous months)
- Simpler implementation
- Prevents UI clutter from too many months
- User can create new months by switching to them

---

## Handling Unsaved Changes

### Confirmation Dialog Pattern

```typescript
function MonthSelector({ currentMonth, onMonthChange }: MonthSelectorProps) {
  const hasUnsavedChanges = useUnsavedChanges()  // Custom hook
  
  const handleChange = (newMonth: string) => {
    if (hasUnsavedChanges) {
      const confirmed = window.confirm(
        'You have unsaved changes. Switch months anyway? Changes will be auto-saved.'
      )
      if (!confirmed) return
    }
    
    onMonthChange(newMonth)
  }
  
  return (
    <select value={currentMonth} onChange={(e) => handleChange(e.target.value)}>
      {/* options */}
    </select>
  )
}

// Custom hook to detect unsaved changes
function useUnsavedChanges(): boolean {
  const [isDirty, setIsDirty] = useState(false)
  
  // Mark dirty when user makes changes
  const markDirty = useCallback(() => setIsDirty(true), [])
  
  // Clear dirty flag when saved
  const markClean = useCallback(() => setIsDirty(false), [])
  
  return isDirty
}
```

**Recommendation**: **Auto-save on switch** instead of confirmation
- Better UX - no annoying dialogs
- Budget apps should auto-save frequently
- Use visual indicator (e.g., "Saving..." text) instead

---

## Visual Feedback During Switch

### Loading States

```typescript
function MonthSelector() {
  const [isLoading, setIsLoading] = useState(false)
  
  const handleMonthChange = async (newMonth: string) => {
    setIsLoading(true)
    
    // Save current month data
    await saveCurrentMonth()
    
    // Load new month data
    await loadMonth(newMonth)
    
    setIsLoading(false)
  }
  
  return (
    <div className="month-selector">
      <select 
        value={currentMonth} 
        onChange={(e) => handleMonthChange(e.target.value)}
        disabled={isLoading}
      >
        {/* options */}
      </select>
      {isLoading && <span className="loading-indicator">Loading...</span>}
    </div>
  )
}
```

### CSS for Loading State

```css
.month-selector {
  display: flex;
  align-items: center;
  gap: 12px;
}

.month-selector select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-indicator {
  font-size: 14px;
  color: #666;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
```

---

## Complete Component Example

```typescript
import { useState, useEffect } from 'react'

interface MonthControlsProps {
  currentMonth: string
  availableMonths: string[]
  onMonthChange: (month: string) => void
}

export function MonthControls({ 
  currentMonth, 
  availableMonths, 
  onMonthChange 
}: MonthControlsProps) {
  const [isLoading, setIsLoading] = useState(false)
  
  const currentIndex = availableMonths.indexOf(currentMonth)
  const hasPrev = currentIndex < availableMonths.length - 1
  const hasNext = currentIndex > 0
  
  const handlePrev = () => {
    if (hasPrev) {
      onMonthChange(availableMonths[currentIndex + 1])
    }
  }
  
  const handleNext = () => {
    if (hasNext) {
      onMonthChange(availableMonths[currentIndex - 1])
    }
  }
  
  const formatMonthDisplay = (monthString: string) => {
    const [year, month] = monthString.split('-')
    const date = new Date(parseInt(year), parseInt(month) - 1)
    return date.toLocaleDateString('en-US', { 
      month: 'long', 
      year: 'numeric' 
    })
  }
  
  return (
    <div className="month-controls">
      <button 
        onClick={handlePrev} 
        disabled={!hasPrev || isLoading}
        aria-label="Previous month"
      >
        ← Prev
      </button>
      
      <select 
        value={currentMonth}
        onChange={(e) => onMonthChange(e.target.value)}
        disabled={isLoading}
        className="month-dropdown"
      >
        {availableMonths.map(month => (
          <option key={month} value={month}>
            {formatMonthDisplay(month)}
          </option>
        ))}
      </select>
      
      <button 
        onClick={handleNext} 
        disabled={!hasNext || isLoading}
        aria-label="Next month"
      >
        Next →
      </button>
      
      {isLoading && (
        <span className="loading-text">Saving...</span>
      )}
    </div>
  )
}
```

---

## Accessibility Considerations

### Keyboard Navigation

- ✅ `<select>` has built-in keyboard support:
  - Arrow keys to navigate options
  - Enter/Space to open dropdown
  - Type-ahead to jump to month starting with letter
- ✅ Prev/Next buttons keyboard accessible with Tab + Enter
- ✅ Use `aria-label` for icon-only buttons

### Screen Reader Support

```typescript
<label id="month-label">Select Budget Month:</label>
<select 
  aria-labelledby="month-label"
  aria-describedby="month-help-text"
  value={currentMonth}
  onChange={handleChange}
>
  {/* options */}
</select>
<span id="month-help-text" className="sr-only">
  Use arrow keys to select a different month. Data is auto-saved when switching.
</span>
```

---

## Pros of This Approach

- ✅ **Simple implementation** - Uses native HTML elements
- ✅ **No dependencies** - No need for React Router or date pickers
- ✅ **Accessible** - Built-in keyboard and screen reader support
- ✅ **Familiar UX** - Users understand dropdowns
- ✅ **Flexible** - Can add prev/next buttons easily
- ✅ **Fast switching** - Direct month selection from dropdown
- ✅ **localStorage integration** - Easy persistence with useEffect
- ✅ **Future-proof** - Can migrate to URL routing later

---

## Cons/Challenges

- ⚠️ **No bookmarkable URLs** - Without routing, can't share specific month URLs
- ⚠️ **No browser back/forward** - Month changes don't update history
- ⚠️ **Limited to ~12 months** - Dropdown becomes unwieldy with many months
- ⚠️ **localStorage limits** - ~5-10MB storage limit per domain

---

## localStorage Persistence Strategy

**See next research doc**: `07-multi-month-persistence.md` for detailed implementation:
- State serialization/deserialization
- Storage quota management
- Conflict resolution
- Migration patterns

---

## Testing Strategy

```typescript
describe('MonthControls', () => {
  it('renders current month in dropdown', () => {
    render(
      <MonthControls 
        currentMonth="2025-10" 
        availableMonths={['2025-10', '2025-09']}
        onMonthChange={vi.fn()}
      />
    )
    expect(screen.getByDisplayValue('October 2025')).toBeInTheDocument()
  })
  
  it('calls onMonthChange when dropdown selection changes', () => {
    const handleChange = vi.fn()
    render(
      <MonthControls 
        currentMonth="2025-10" 
        availableMonths={['2025-10', '2025-09']}
        onMonthChange={handleChange}
      />
    )
    
    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: '2025-09' } })
    
    expect(handleChange).toHaveBeenCalledWith('2025-09')
  })
  
  it('navigates to previous month with prev button', () => {
    const handleChange = vi.fn()
    render(
      <MonthControls 
        currentMonth="2025-10" 
        availableMonths={['2025-10', '2025-09', '2025-08']}
        onMonthChange={handleChange}
      />
    )
    
    fireEvent.click(screen.getByLabelText('Previous month'))
    expect(handleChange).toHaveBeenCalledWith('2025-09')
  })
  
  it('disables prev button when at earliest month', () => {
    render(
      <MonthControls 
        currentMonth="2025-08" 
        availableMonths={['2025-10', '2025-09', '2025-08']}
        onMonthChange={vi.fn()}
      />
    )
    
    expect(screen.getByLabelText('Previous month')).toBeDisabled()
  })
})
```

---

## Next Steps

1. **Implement MonthControls component** with dropdown + prev/next buttons
2. **Add month state to app root** with `useState`
3. **Implement auto-save on switch** using `useEffect`
4. **Research localStorage persistence** (next doc: 07-multi-month-persistence.md)
5. **Add visual loading feedback** during month switch
6. **Test with multiple months** to ensure data isolation

---

## References

- [React: Choosing State Structure](https://react.dev/learn/choosing-the-state-structure)
- [React: Synchronizing with Effects](https://react.dev/learn/synchronizing-with-effects)
- [MDN: localStorage](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage)
- [React: `<select>` Reference](https://react.dev/reference/react-dom/components/select)
