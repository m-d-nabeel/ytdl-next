"use client";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useRouter } from "next/navigation";

const formSchema = z.object({
  url: z.string().url(),
});
export default function LinkForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      url: "",
    },
  });

  const router = useRouter();

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      const response = await fetch("/api/download", {
        method: "POST",
        body: JSON.stringify({
          url: values.url,
        }),
      });
      if (!response.ok) {
        console.log("Illegal Response from Server");
        return;
      }
      const data = await response.json();
      let { title, info } = data;
      if (!title) {
        console.log("No title in Response");
        return;
      }
      title = encodeURIComponent(title.replaceAll("_", "="));
      router.push(`/video/${title}`);
    } catch (error) {
      console.log("Form Submit Error");
    }
  }

  if (form.formState.isSubmitting) {
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <div className="w-32 h-32 border-8 border-gray-300 border-t-slate-800 rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 w-full">
        <FormField
          control={form.control}
          name="url"
          render={({ field }) => (
            <FormItem>
              <FormLabel>URL</FormLabel>
              <FormControl>
                <Input placeholder="Enter url" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}
