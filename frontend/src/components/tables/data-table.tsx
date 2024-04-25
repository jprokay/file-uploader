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
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { useQuery, useQueryClient } from "@tanstack/react-query"

interface DataTableProps<TData, TValue> {
	columns: ColumnDef<TData, TValue>[]
	data: TData[]
	queryFn: (state: PaginationState) => () => Promise<TData[]>
	pageSizes: Array<number>
	defaultPageSize: number
	queryKey: string
	rowCount: number
}

export function DataTable<TData, TValue>({
	queryKey,
	columns,
	data,
	queryFn,
	pageSizes,
	defaultPageSize,
	rowCount
}: DataTableProps<TData, TValue>) {
	const [pagination, setPagination] = useState<PaginationState>({
		pageIndex: 0,
		pageSize: defaultPageSize
	})

	const dataQuery = useQuery({
		queryKey: [queryKey, pagination],
		queryFn: queryFn(pagination),
		initialData: data,
	})

	const table = useReactTable({
		data: dataQuery.data,
		columns,
		getCoreRowModel: getCoreRowModel(),
		manualPagination: true,
		debugTable: true,
		state: {
			pagination,
		},
		onPaginationChange: setPagination,
		rowCount
	})

	const queryClient = useQueryClient()

	const prefetch = (newPagination: PaginationState) => {
		queryClient.prefetchQuery({
			queryKey: [queryKey, newPagination],
			queryFn: queryFn(newPagination),
			// Prefetch only fires when data is older than the staleTime,
			// so in a case like this you definitely want to set one
			staleTime: 60000,
		})
	}

	const prefetchNextPage = () => {
		return prefetch({ ...pagination, pageIndex: pagination.pageIndex + 1 })
	}

	return (
		<div>
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
				<div className="flex items-center justify-center space-x-2 py-4 border-t-black">
					<Button
						variant="outline"
						size="sm"
						onClick={() => table.previousPage()}
						disabled={!table.getCanPreviousPage()}
					>
						Previous
					</Button>
					<Select onValueChange={(v) => setPagination((prev) => ({ ...prev, pageSize: Number(v) }))} value={pagination.pageSize.toString()}>
						<SelectTrigger className="w-[180px]">
							<SelectValue placeholder="Page Size" />
						</SelectTrigger>
						<SelectContent>
							<SelectGroup>
								<SelectLabel>Page Size</SelectLabel>
								{pageSizes.map((size) => <SelectItem key={size} value={size.toString()}>{size} Rows</SelectItem>)}
							</SelectGroup>
						</SelectContent>
					</Select>
					<Button
						variant="outline"
						size="sm"
						onClick={() => table.nextPage()}
						disabled={!table.getCanNextPage()}
						onMouseEnter={prefetchNextPage}
						onFocus={prefetchNextPage}
					>
						Next
					</Button>
				</div>

			</div>
		</div>
	)
}

