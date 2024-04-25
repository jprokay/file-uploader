"use client"

import {
	ColumnDef,
	flexRender,
	getCoreRowModel,
	useReactTable,
	PaginationState
} from "@tanstack/react-table"

import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { useMemo, useState } from "react"
import { keepPreviousData, useQuery, useQueryClient } from "@tanstack/react-query"
import { ChevronLeftIcon, ChevronRightIcon, ChevronsLeft, ChevronsRight } from "lucide-react"
import { Input } from "../ui/input"
import { Label } from "../ui/label"

interface ItemsWithTotal<TData> {
	items: TData[]
	total: number
}

interface DataTableProps<TData, TValue> {
	columns: ColumnDef<TData, TValue>[]
	queryFn: (state: PaginationState, search: string | undefined) => () => Promise<ItemsWithTotal<TData>>
	pageSizes: Array<number>
	defaultPageSize: number
	queryKey: string
	enableSearch: boolean
}

export function DataTable<TData, TValue>({
	queryKey,
	columns,
	queryFn,
	pageSizes,
	defaultPageSize,
	enableSearch
}: DataTableProps<TData, TValue>) {
	const [pagination, setPagination] = useState<PaginationState>({
		pageIndex: 0,
		pageSize: defaultPageSize
	})

	const [searchQuery, setSearchQuery] = useState<string>()

	const dataQuery = useQuery({
		queryKey: [queryKey, pagination, searchQuery],
		queryFn: queryFn(pagination, searchQuery),
		placeholderData: keepPreviousData
	})
	const defaultData = useMemo(() => [], [])

	const table = useReactTable({
		data: dataQuery?.data?.items || defaultData,
		columns,
		getCoreRowModel: getCoreRowModel(),
		manualPagination: true,
		debugTable: true,
		state: {
			pagination,
		},
		onPaginationChange: setPagination,
		rowCount: dataQuery?.data?.total || 0,
	})

	const queryClient = useQueryClient()

	const prefetch = (newPagination: PaginationState) => {
		queryClient.prefetchQuery({
			queryKey: [queryKey, newPagination],
			queryFn: queryFn(newPagination, searchQuery),
			// Prefetch only fires when data is older than the staleTime,
			// so in a case like this you definitely want to set one
			staleTime: 60_000,
		})
	}

	const prefetchNextPage = () => {
		return prefetch({ ...pagination, pageIndex: pagination.pageIndex + 1 })
	}

	return (
		<div>
			{enableSearch &&
				(
					<div className="grid w-full max-w-sm items-center gap-1.5 my-4">
						<Label htmlFor="search">Filter Contacts</Label>
						<div className="flex w-full max-w-sm items-center space-x-2">
							<Input id="search" type="" placeholder="Search" onChange={(e) => setSearchQuery(e.target.value)} />
							<Button type="button" onClick={() => setSearchQuery(undefined)}>Reset</Button>
						</div>
					</div>
				)}
			<div className="rounded-md border">
				<Table>
					<TableHeader>
						{table.getHeaderGroups().map((headerGroup) => (
							<TableRow key={headerGroup.id}>
								{headerGroup.headers.map((header) => {
									return (
										<TableHead key={header.id}>
											{header.isPlaceholder
												? null
												: flexRender(
													header.column.columnDef.header,
													header.getContext()
												)}
										</TableHead>
									)
								})}
							</TableRow>
						))}
					</TableHeader>
					<TableBody>
						{table.getRowModel().rows?.length ? (
							table.getRowModel().rows.map((row) => (
								<TableRow
									key={row.id}
									data-state={row.getIsSelected() && "selected"}
									className="even:bg-muted"
								>
									{row.getVisibleCells().map((cell) => (
										<TableCell key={cell.id}>
											{flexRender(cell.column.columnDef.cell, cell.getContext())}
										</TableCell>
									))}
								</TableRow>
							))
						) : (
							<TableRow>
								<TableCell colSpan={columns.length} className="h-24 text-center">
									No results.
								</TableCell>
							</TableRow>
						)}
					</TableBody>
				</Table>
				<div className="flex items-center justify-evenly px-2 my-4">
					<div className="flex items-center space-x-6 lg:space-x-8">
						<div className="flex items-center space-x-2">
							<p className="text-sm font-medium">Rows per page</p>
							<Select
								value={`${table.getState().pagination.pageSize}`}
								onValueChange={(value) => {
									table.setPageSize(Number(value))
								}}
							>
								<SelectTrigger className="h-8 w-[70px]">
									<SelectValue placeholder={table.getState().pagination.pageSize} />
								</SelectTrigger>
								<SelectContent side="top">
									{pageSizes.map((pageSize) => (
										<SelectItem key={pageSize} value={`${pageSize}`}>
											{pageSize}
										</SelectItem>
									))}
								</SelectContent>
							</Select>
						</div>
						<div className="flex w-[100px] items-center justify-center text-sm font-medium">
							Page {table.getState().pagination.pageIndex + 1} of{" "}
							{table.getPageCount()}
						</div>
						<div className="flex items-center space-x-2">
							<Button
								variant="outline"
								className="hidden h-8 w-8 p-0 lg:flex"
								onClick={() => table.setPageIndex(0)}
								disabled={!table.getCanPreviousPage()}
							>
								<span className="sr-only">Go to first page</span>
								<ChevronsLeft className="h-4 w-4" />
							</Button>
							<Button
								variant="outline"
								className="h-8 w-8 p-0"
								onClick={() => table.previousPage()}
								disabled={!table.getCanPreviousPage()}
							>
								<span className="sr-only">Go to previous page</span>
								<ChevronLeftIcon className="h-4 w-4" />
							</Button>
							<Button
								variant="outline"
								className="h-8 w-8 p-0"
								onClick={() => table.nextPage()}
								disabled={!table.getCanNextPage()}
								onMouseEnter={prefetchNextPage}
								onFocus={prefetchNextPage}
							>
								<span className="sr-only">Go to next page</span>
								<ChevronRightIcon className="h-4 w-4" />
							</Button>
							<Button
								variant="outline"
								className="hidden h-8 w-8 p-0 lg:flex"
								onClick={() => table.setPageIndex(table.getPageCount() - 1)}
								disabled={!table.getCanNextPage()}
							>
								<span className="sr-only">Go to last page</span>

								<ChevronsRight className="h-4 w-4" />
							</Button>
						</div>
					</div>
				</div>
			</div>
		</div>
	)
}

