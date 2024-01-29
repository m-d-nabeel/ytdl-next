"use client";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { set, useForm } from "react-hook-form";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useStore } from "@/hooks/use-store";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { YTVideoDetail } from "@/custom-types";

const formSchema = z.object({
  url: z.string().url(),
  quality: z.enum(["low", "medium", "high", "audio_only"]),
});
export default function LinkForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      url: "",
      quality: "medium",
    },
  });

  const [fetched, setIsFetched] = useState(false);
  const router = useRouter();

  const { setVideoDetails, videoDetails } = useStore();

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      const response = await fetch(
        "/api/download/" + encodeURIComponent(values.url),
        {
          method: "POST",
          body: JSON.stringify({ quality: values.quality }),
        }
      );
      if (!response.ok) {
        console.log("Illegal Response(Not Ok) from Server");
        return;
      }
      response
        .json()
        .then((data: YTVideoDetail) => {
          if (!data.title) {
            console.log("No title in Response");
            return;
          }
          setVideoDetails(data);
          setIsFetched(true);
        })
        .catch((error) => {
          console.log("Error Parsing Response");
        });
    } catch (error) {
      console.log("Form Submit Error");
    }
  }

  if (fetched) {
    // videoDetails can't be null here, cause if fetched is true,
    // it means that videoDetails is set
    router.push(`/video/${encodeURIComponent(videoDetails!.title)}`);
    return null;
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
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="gap-3 w-full grid grid-cols-2"
      >
        <FormField
          control={form.control}
          name="url"
          render={({ field }) => (
            <FormItem className="col-span-2">
              <FormLabel className="text-sm font-semibold">URL</FormLabel>
              <FormControl>
                <Input placeholder="Enter url" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="quality"
          render={({ field }) => (
            <FormItem className="col-span-2">
              <FormLabel className="text-sm font-semibold">QUALITY</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl className="uppercase font-semibold text-sm">
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent defaultValue={"MEDIUM"}>
                  {formSchema.shape.quality.options.map((option) => (
                    <SelectItem
                      className="uppercase font-semibold text-sm"
                      key={option}
                      value={option}
                    >
                      {option}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}
