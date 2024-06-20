CREATE TYPE role AS ENUM ('user', 'admin', 'organizer');

CREATE TABLE IF NOT EXISTS public.users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  username VARCHAR(20) NOT NULL UNIQUE,
  email VARCHAR(256) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  birth_date DATE,
  role role NOT NULL DEFAULT 'user',
  verified boolean DEFAULT 'false',
  about VARCHAR(500),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE event_type AS ENUM ('offline', 'online', 'both');

CREATE TABLE IF NOT EXISTS public.events (
   id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   name VARCHAR(500) NOT NULL,
   organizer_id UUID NOT NULL,
   description VARCHAR(1000),
   start_date TIMESTAMPTZ NOT NULL,
   end_date TIMESTAMPTZ NOT NULL,
   is_paid BOOLEAN default 'false',
   event_type event_type NOT NULL DEFAULT 'online',
   country varchar(5),
   city varchar(50), 
   slug text not null,
   likes INT NOT NULL DEFAULT 0,
   follows INT NOT NULL DEFAULT 0,
   attendees INT NOT NULL DEFAULT 0,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (organizer_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.reviews (
   id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   title VARCHAR(500) NOT NULL,
   event_id UUID NOT NULL,
   author_id UUID NOT NULL,
   body VARCHAR(2000) NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (author_id) REFERENCES public.users(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.event_attendees (
   event_id UUID NOT NULL,
   attendee_id UUID NOT NULL,
   FOREIGN KEY (attendee_id) REFERENCES public.users(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE,
   PRIMARY KEY(event_id, attendee_id),
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.event_followers (
   event_id UUID NOT NULL,
   follower_id UUID NOT NULL,
   FOREIGN KEY (follower_id) REFERENCES public.users(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE,
   PRIMARY KEY(event_id, follower_id)
);

CREATE TABLE IF NOT EXISTS public.event_likes (
   event_id UUID NOT NULL,
   user_id UUID NOT NULL,
   FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE,
   PRIMARY KEY(event_id, user_id)
);

CREATE TABLE IF NOT EXISTS public.tags (
   id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.categories (
   id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.event_tags (
   event_id UUID NOT NULL,
   tag_id UUID NOT NULL,
   FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE,
   PRIMARY KEY(event_id, tag_id)
);

CREATE TABLE IF NOT EXISTS public.event_categories (
   event_id UUID NOT NULL,
   category_id UUID NOT NULL,
   FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE,
   FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE,
   PRIMARY KEY(event_id, category_id)
);

-- ============================================================================================================
-- Function & Trigger to update the likes count when LIKE is added / removed.
-- ============================================================================================================
-- > Function
CREATE OR REPLACE FUNCTION update_event_likes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE public.events
        SET likes = likes + 1
        WHERE id = NEW.event_id;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE public.events
        SET likes = GREATEST(likes - 1, 0)
        WHERE id = OLD.event_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================================================
-- > Trigger
CREATE TRIGGER update_likes_count
AFTER INSERT OR DELETE ON public.event_likes
FOR EACH ROW
EXECUTE FUNCTION update_event_likes_count();

-- ============================================================================================================
-- Function & Trigger to update the follows count when FOLLOW is added / removed.
-- ============================================================================================================
-- > Function
CREATE OR REPLACE FUNCTION update_event_follows_count()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE public.events
        SET follows = follows + 1
        WHERE id = NEW.event_id;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE public.events
        SET follows = GREATEST(follows - 1, 0)
        WHERE id = OLD.event_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================================================
-- > Trigger
CREATE TRIGGER update_follows_count
AFTER INSERT OR DELETE ON public.event_followers
FOR EACH ROW
EXECUTE FUNCTION update_event_follows_count();

-- ============================================================================================================
-- Function & Trigger to update the attendees count when attendee is added / removed.
-- ============================================================================================================
-- > Function
CREATE OR REPLACE FUNCTION update_event_attendees_count()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE public.events
        SET attendees = attendees + 1
        WHERE id = NEW.event_id;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE public.events
        SET attendees = GREATEST(attendees - 1, 0)
        WHERE id = OLD.event_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================================================
-- > Trigger
CREATE TRIGGER update_attendees_count
AFTER INSERT OR DELETE ON public.event_attendees
FOR EACH ROW
EXECUTE FUNCTION update_event_attendees_count();

-- ============================================================================================================