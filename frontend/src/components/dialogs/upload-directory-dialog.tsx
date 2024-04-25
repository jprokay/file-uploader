import { Button } from "@/components/ui/button"
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog"
import { UploadDirectoryForm } from "../forms/upload-directory-form"
import { UploadCloud } from "lucide-react"

export function UploadDirectoryDemo() {
	return (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline" className="border-zinc-600 text-zinc-600">
					<UploadCloud className="w-4 h-4 mr-4" />Upload</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px] md:max-w-[650px] lg:max-w-[850px]">
				<DialogHeader>
					<DialogTitle>Upload New Directory</DialogTitle>
					<DialogDescription>
						Select a CSV to upload. Columns are expected to be in the order: first_name, last_name, email
					</DialogDescription>
				</DialogHeader>
				<div className="grid gap-4 py-4">
					<UploadDirectoryForm />
				</div>
			</DialogContent>
		</Dialog>
	)
}

