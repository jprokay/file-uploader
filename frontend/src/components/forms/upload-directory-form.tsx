"use client"
import { Button } from "@/components/ui/button"
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { client } from "@/lib/contacts-api"
import { getCookieValue } from "@/lib/cookies"
import { useQueryClient } from "@tanstack/react-query"
import { useState } from "react"
import { CsvPreviewTable } from "../tables/csv-preview-table"
import { Checkbox } from "../ui/checkbox"

const formSchema = z.object({
	fileNames: z
		.custom<FileList>((v) => v instanceof FileList, {
			message: 'Files are required',
		}),
	excludeFirstRow: z.boolean(),
	columnLayout: z.array(z.string())
})

export function UploadDirectoryForm() {

	const [preview, setPreview] = useState<Array<Array<string>>>([]);

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			fileNames: undefined,
			excludeFirstRow: true,
			columnLayout: ["first_name", "last_name", "email"]
		}
	})

	const queryClient = useQueryClient()

	async function onSubmit(values: z.infer<typeof formSchema>) {
		const files: Array<File> = []
		for (let i = 0; i < values.fileNames.length; i++) {
			const item = values.fileNames.item(i)
			if (item !== null) {
				files.push(item)
			}
		}

		await client.POST("/directories", {
			body: {
				columnLayout: ["first_name", "last_name", "email"],
				excludeFirstRow: values.excludeFirstRow,
				filename: files.map((f) => f.name)
			}, params: {
				cookie: { userId: getCookieValue("userId") || "" }
			},
			bodySerializer(_body) {
				const fd = new FormData();
				for (let i = 0; i < values.fileNames.length; i++) {
					const item = values.fileNames.item(i)
					if (item !== null) {
						fd.append("filename", item)
					}
				}

				fd.set("excludeFirstRow", String(values.excludeFirstRow))
				fd.set("columnLayout", new Blob(values.columnLayout, { type: "text/plain" }));
				return fd;
			},
		})
		await queryClient.invalidateQueries({ queryKey: ['directories'] })
	}

	const excludeFirstRow = form.watch("excludeFirstRow", true)

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
				<FormField
					control={form.control}
					name="fileNames"
					render={({ field: { ref, name, onBlur, onChange } }) => (
						<FormItem>
							<FormLabel>Files</FormLabel>
							<FormControl>
								<Input type="file"
									ref={ref}
									name={name}
									onBlur={onBlur}
									multiple={false}
									onChange={(e) => {
										onChange(e.target.files)
										if (e.target.files) {
											const file = e.target.files[0]
											const reader = new FileReader();

											reader.onload = function(e) {
												const contents: string = String(e!.target!.result)
												const rows = contents.split('\n');
												const preview = rows.slice(0, 5);
												setPreview(preview.map((row) => row.split(',')));
											}
											reader.readAsText(file)
										}
									}}
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="excludeFirstRow"
					render={({ field }) => (
						<FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
							<FormLabel>Exclude First Row?</FormLabel>
							<FormControl>
								<Checkbox checked={field.value} onCheckedChange={field.onChange} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<CsvPreviewTable rows={preview} excludeFirstRow={excludeFirstRow} />
				<Button type="submit" disabled={form.formState.isSubmitting}>Submit</Button>
			</form>
		</Form>
	)
}
