
-- remove triggers and functions

DROP TRIGGER IF EXISTS update_attendees_count ON public.event_attendees;
DROP FUNCTION IF EXISTS update_event_attendees_count;
DROP TRIGGER IF EXISTS update_follows_count ON public.event_followers;
DROP FUNCTION IF EXISTS update_event_follows_count;
DROP TRIGGER IF EXISTS update_likes_count ON public.event_likes;
DROP FUNCTION IF EXISTS update_event_likes_count;

-- remove tables

DROP TABLE IF EXISTS public.event_categories;
DROP TABLE IF EXISTS public.event_tags;
DROP TABLE IF EXISTS public.categories;
DROP TABLE IF EXISTS public.tags;
DROP TABLE IF EXISTS public.event_likes;
DROP TABLE IF EXISTS public.event_attendees;
DROP TABLE IF EXISTS public.event_followers;
DROP TABLE IF EXISTS public.reviews;
DROP TABLE IF EXISTS public.events;
DROP TABLE IF EXISTS public.reviews;
DROP TABLE IF EXISTS public.users;

-- remove enums 

DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS role;
